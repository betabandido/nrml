package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func HealthCheck() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if ctx.Request.URL.Path == "/health" {
			ctx.String(http.StatusOK, "Healthy")
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}
