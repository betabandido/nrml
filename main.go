package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"net/http"
	"nrml/api"
	"nrml/app"
	"nrml/config"
	"nrml/logging"
	"os"
	"time"
)

func init() {
	err := config.ReadConfiguration(app.Name)
	if err != nil {
		_, _ = fmt.Fprint(os.Stderr, err.Error())
		os.Exit(1)
	}

	err = logging.SetupWithOptions(
		logging.Options{
			Labels: viper.GetStringMapString("newRelic.labels"),
		},
	)
	if err != nil {
		_, _ = fmt.Fprint(os.Stderr, err.Error())
		os.Exit(1)
	}
}

func main() {
	log.Infof("Starting service")

	tableName := viper.GetString("tableName")
	if tableName == "" {
		log.Fatal("no table name was specified")
	}

	awsRegion := viper.GetString("aws.region")
	if awsRegion == "" {
		log.Fatal("no AWS region was specified")
	}

	theApp := app.New(tableName, awsRegion)
	engine := gin.New()

	api.Setup(engine, theApp)

	server := &http.Server{
		Addr:         ":8000",
		Handler:      engine,
		ReadTimeout:  2 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	log.Fatal(server.ListenAndServe())
}
