package middleware

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func CorsConfiguration(allowedHeaders []string) gin.HandlerFunc {
	config := cors.DefaultConfig()
	config.AllowHeaders = append(config.AllowHeaders, allowedHeaders...)
	config.AllowAllOrigins = true

	return cors.New(config)
}
