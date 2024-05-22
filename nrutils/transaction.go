package nrutils

import (
	"github.com/gin-gonic/gin"
	"github.com/newrelic/go-agent/v3/integrations/nrgin"
)

type AnnotationConfig struct {
	QueryParameters []string
	HeaderNames     []string
}

func AnnotateTransactions(annotations AnnotationConfig) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		txn := nrgin.Transaction(ctx)

		for _, queryParamName := range annotations.QueryParameters {
			if queryParamValue := ctx.Query(queryParamName); queryParamValue != "" {
				txn.AddAttribute(queryParamName, queryParamValue)
			}
		}

		for _, headerName := range annotations.HeaderNames {
			if headerValue := ctx.GetHeader(headerName); headerValue != "" {
				txn.AddAttribute(headerName, headerValue)
			}
		}

		ctx.Next()

		for _, err := range ctx.Errors {
			txn.NoticeError(err)
		}
	}
}
