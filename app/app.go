package app

import (
	"github.com/newrelic/go-agent/v3/newrelic"
	"nrml/nrutils"
	"nrml/products"
)

const Name = "NRML"

type App struct {
	ProductGetter       products.ProductGetter
	NewRelicApplication *newrelic.Application
}

func New(tableName string, awsRegion string) *App {
	return &App{
		ProductGetter:       products.NewDefaultRepository(tableName, awsRegion),
		NewRelicApplication: nrutils.CreateNewRelicApplication(),
	}
}
