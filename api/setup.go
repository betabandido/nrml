package api

import (
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"nrml/app"
	"nrml/middleware"
	"nrml/nrutils"
	"os"
)

func Setup(
	engine *gin.Engine,
	theApp *app.App,
) {
	engine.Use(gin.Recovery())
	engine.Use(middleware.HealthCheck())
	if os.Getenv("NRML_ENABLE_PPROF") == "true" {
		pprof.Register(engine)
	}

	nrutils.SetupGinEngine(
		engine,
		theApp.NewRelicApplication,
		[]string{},
	)

	engine.Use(middleware.CorsConfiguration([]string{"newrelic"}))
	engine.Use(middleware.LoggingMiddleware())

	v1 := engine.Group("/api/v1")

	NewProductHandler(theApp.ProductGetter).
		RegisterRoutes(v1)
}
