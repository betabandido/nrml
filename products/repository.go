//go:generate mockery --all --output ./mocks
package products

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	log "github.com/sirupsen/logrus"
	"nrml/nrutils"
	"strings"
)

const (
	productKeyPartitionKey = "ProductKey"
)

type ProductGetter interface {
	GetProductByProductKey(
		ctx context.Context,
		tenant string,
		locale string,
		productKey string,
	) (*Product, error)
}

type DefaultRepository struct {
	dbClient  *dynamodb.Client
	tableName string
}

func NewDefaultRepository(tableName string, awsRegion string) *DefaultRepository {
	dbClient := createDynamoDBClient(awsRegion)

	return &DefaultRepository{
		dbClient:  dbClient,
		tableName: tableName,
	}
}

func (r *DefaultRepository) GetProductByProductKey(
	ctx context.Context,
	tenant string,
	locale string,
	productKey string,
) (*Product, error) {
	expr, err := buildExpression(tenant, locale, productKey, productKeyPartitionKey)
	if err != nil {
		return nil, err
	}

	output, err := r.dbClient.Query(ctx, &dynamodb.QueryInput{
		TableName:                 aws.String(r.tableName),
		KeyConditionExpression:    expr.KeyCondition(),
		FilterExpression:          expr.Filter(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		ScanIndexForward:          aws.Bool(true),
	})
	if err != nil {
		return nil, err
	}

	if output.Items == nil || len(output.Items) == 0 {
		return nil, nil
	}

	var product ProductDbItem
	err = attributevalue.UnmarshalMap(output.Items[0], &product)
	if err != nil {
		return nil, err
	}

	return &product.Product, nil
}

func createDynamoDBClient(awsRegion string) *dynamodb.Client {
	cfg, err := config.LoadDefaultConfig(
		context.Background(),
		config.WithRegion(awsRegion),
	)
	if err != nil {
		log.Fatalf("Error configuring AWS: %s", err)
	}

	nrutils.InstrumentAWS(&cfg)

	return dynamodb.NewFromConfig(cfg)
}

func buildKey(tenant string, locale string, key string) (string, error) {
	if tenant == "" {
		return "", fmt.Errorf("invalid tenant: %s", tenant)
	}
	if locale == "" {
		return "", fmt.Errorf("invalid locale: %s", locale)
	}
	if key == "" {
		return "", fmt.Errorf("invalid key: %s", key)
	}
	return fmt.Sprintf("%s#%s#%s", strings.ToLower(tenant), strings.ToLower(locale), key), nil
}

func buildExpression(
	tenant, locale, key, partitionKey string,
) (expression.Expression, error) {
	dbKey, err := buildKey(tenant, locale, key)
	if err != nil {
		return expression.Expression{}, err
	}

	builder := expression.NewBuilder().WithKeyCondition(
		expression.Key(partitionKey).
			Equal(expression.Value(dbKey)),
	)

	return builder.Build()
}
