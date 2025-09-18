package middleware

import (
	"time"
	"user/infrastructure/log"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
)

func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := uuid.New().String()
		timeoutCtx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
		defer cancel()

		ctx := context.WithValue(timeoutCtx, "requestID", requestID)
		c.Request = c.Request.WithContext(ctx)

		startTime := time.Now()
		c.Next()
		latency := time.Since(startTime)
		requestLog := logrus.Fields{
			"request_id": requestID,
			"method":     c.Request.Method,
			"path":       c.Request.URL.Path,
			"status":     c.Writer.Status(),
			"latency":    latency,
			"client_ip":  c.ClientIP(),
		}

		if c.Writer.Status() == 200 || c.Writer.Status() == 201 {
			log.Logger.WithFields(requestLog).Info("Request completed")
		} else {
			log.Logger.WithFields(requestLog).Info("Request failed")
		}
	}
}
