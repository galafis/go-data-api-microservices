package middleware

import (
	"time"

	"github.com/galafis/go-data-api-microservices/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// RequestID is a middleware that adds a request ID to the context
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check if request ID is already set
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String()
		}

		// Set request ID in context and header
		c.Set("request_id", requestID)
		c.Header("X-Request-ID", requestID)

		c.Next()
	}
}

// Logger is a middleware that logs request information
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start timer
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// Get request ID
		requestID, exists := c.Get("request_id")
		if !exists {
			requestID = uuid.New().String()
			c.Set("request_id", requestID)
		}

		// Process request
		c.Next()

		// Calculate latency
		latency := time.Since(start)

		// Get client IP
		clientIP := c.ClientIP()

		// Get status code
		statusCode := c.Writer.Status()

		// Get error if any
		var errorMessage string
		if len(c.Errors) > 0 {
			errorMessage = c.Errors.String()
		}

		// Log request
		logger.WithFields(map[string]interface{}{
			"request_id":  requestID,
			"status_code": statusCode,
			"latency":     latency,
			"client_ip":   clientIP,
			"method":      c.Request.Method,
			"path":        path,
			"query":       raw,
			"user_agent":  c.Request.UserAgent(),
			"error":       errorMessage,
		}).Info("Request processed")
	}
}

// Recovery is a middleware that recovers from panics
func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Get request ID
				requestID, exists := c.Get("request_id")
				if !exists {
					requestID = uuid.New().String()
				}

				// Log error
				logger.WithFields(map[string]interface{}{
					"request_id": requestID,
					"error":      err,
					"path":       c.Request.URL.Path,
					"method":     c.Request.Method,
					"client_ip":  c.ClientIP(),
				}).Error("Panic recovered")

				// Return error to client
				c.AbortWithStatusJSON(500, gin.H{
					"error":      "Internal server error",
					"request_id": requestID,
				})
			}
		}()

		c.Next()
	}
}

