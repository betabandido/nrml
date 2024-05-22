package nrutils

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/gin-gonic/gin"
	nrlogrusv2 "github.com/newrelic/go-agent/v3/integrations/logcontext-v2/nrlogrus"
	"github.com/newrelic/go-agent/v3/integrations/nrawssdk-v2"
	"github.com/newrelic/go-agent/v3/integrations/nrgin"
	"github.com/newrelic/go-agent/v3/newrelic"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func CreateNewRelicApplication() *newrelic.Application {
	newRelicApp, err := newrelic.NewApplication(
		newrelic.ConfigAppName(viper.GetString("newRelic.appName")),
		newrelic.ConfigLicense(viper.GetString("newRelic.license")),
		newrelic.ConfigDistributedTracerEnabled(true),
		newrelic.ConfigAppLogDecoratingEnabled(true),
		newrelic.ConfigAppLogForwardingEnabled(false),
		func(config *newrelic.Config) {
			config.Labels = viper.GetStringMapString("newRelic.labels")
		},
	)
	if err != nil {
		log.Fatalf("Error setting up New Relic: %s", err)
	}
	log.SetFormatter(nrlogrusv2.NewFormatter(newRelicApp, &log.JSONFormatter{}))

	return newRelicApp
}

func SetupGinEngine(
	engine *gin.Engine,
	newRelicApp *newrelic.Application,
	queryParameters []string,
) {
	annotations := AnnotationConfig{
		QueryParameters: queryParameters,
		HeaderNames:     nil,
	}
	SetupGinEngineWithAnnotations(engine, newRelicApp, annotations)
}

func SetupGinEngineWithAnnotations(
	engine *gin.Engine,
	newRelicApp *newrelic.Application,
	annotations AnnotationConfig,
) {
	engine.Use(nrgin.Middleware(newRelicApp))
	engine.Use(AnnotateTransactions(annotations))
}

func InstrumentAWS(awsConfig *aws.Config) {
	nrawssdk.AppendMiddlewares(&awsConfig.APIOptions, nil)
}
