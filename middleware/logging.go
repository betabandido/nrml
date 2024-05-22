package middleware

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"time"
)

func LoggingMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		t := time.Now()

		ctx.Next()

		latency := time.Since(t)

		entry := log.WithContext(ctx).
			WithFields(log.Fields{
				"method":     ctx.Request.Method,
				"path":       ctx.Request.URL.Path,
				"query":      ctx.Request.URL.RawQuery,
				"status":     ctx.Writer.Status(),
				"latency_us": latency.Nanoseconds() / 1000,
			})

		if ctx.Writer.Status() < 400 {
			entry.Info("Request completed")
		} else {
			if len(ctx.Errors) > 0 {
				entry = entry.WithField("errors", ctx.Errors)
			}
			entry.Error("Error executing request")
		}
	}
}
