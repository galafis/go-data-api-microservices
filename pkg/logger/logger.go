package logger

import (
	"context"
	"fmt"
	"io"
	"os"
	"runtime"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

var (
	// Logger is the global logger instance
	logger *logrus.Logger
)

// Config represents logger configuration
type Config struct {
	Level      string
	Format     string
	Output     string
	TimeFormat string
}

// init initializes the logger with default configuration
func init() {
	logger = logrus.New()
	logger.SetLevel(logrus.InfoLevel)
	logger.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: time.RFC3339,
	})
	logger.SetOutput(os.Stdout)
}

// Configure configures the logger with the provided configuration
func Configure(config *Config) {
	// Set log level
	level, err := logrus.ParseLevel(config.Level)
	if err != nil {
		level = logrus.InfoLevel
	}
	logger.SetLevel(level)

	// Set formatter
	switch config.Format {
	case "json":
		logger.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: config.TimeFormat,
		})
	case "text":
		logger.SetFormatter(&logrus.TextFormatter{
			TimestampFormat: config.TimeFormat,
			FullTimestamp:   true,
		})
	default:
		logger.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: config.TimeFormat,
		})
	}

	// Set output
	switch config.Output {
	case "stdout":
		logger.SetOutput(os.Stdout)
	case "stderr":
		logger.SetOutput(os.Stderr)
	default:
		// Try to open file
		file, err := os.OpenFile(config.Output, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			logger.SetOutput(os.Stdout)
			Error("Failed to open log file, using stdout instead", err)
		} else {
			logger.SetOutput(file)
		}
	}
}

// SetOutput sets the logger output
func SetOutput(output io.Writer) {
	logger.SetOutput(output)
}

// SetLevel sets the logger level
func SetLevel(level string) {
	lvl, err := logrus.ParseLevel(level)
	if err != nil {
		lvl = logrus.InfoLevel
	}
	logger.SetLevel(lvl)
}

// WithFields returns a new entry with the specified fields
func WithFields(fields map[string]interface{}) *logrus.Entry {
	return logger.WithFields(logrus.Fields(fields))
}

// WithContext returns a new entry with the specified context
func WithContext(ctx context.Context) *logrus.Entry {
	return logger.WithContext(ctx)
}

// Debug logs a debug message
func Debug(args ...interface{}) {
	logger.Debug(args...)
}

// Debugf logs a formatted debug message
func Debugf(format string, args ...interface{}) {
	logger.Debugf(format, args...)
}

// Info logs an info message
func Info(args ...interface{}) {
	logger.Info(args...)
}

// Infof logs a formatted info message
func Infof(format string, args ...interface{}) {
	logger.Infof(format, args...)
}

// Warn logs a warning message
func Warn(args ...interface{}) {
	logger.Warn(args...)
}

// Warnf logs a formatted warning message
func Warnf(format string, args ...interface{}) {
	logger.Warnf(format, args...)
}

// Error logs an error message
func Error(args ...interface{}) {
	logger.Error(args...)
}

// Errorf logs a formatted error message
func Errorf(format string, args ...interface{}) {
	logger.Errorf(format, args...)
}

// Fatal logs a fatal message and exits
func Fatal(args ...interface{}) {
	logger.Fatal(args...)
}

// Fatalf logs a formatted fatal message and exits
func Fatalf(format string, args ...interface{}) {
	logger.Fatalf(format, args...)
}

// Panic logs a panic message and panics
func Panic(args ...interface{}) {
	logger.Panic(args...)
}

// Panicf logs a formatted panic message and panics
func Panicf(format string, args ...interface{}) {
	logger.Panicf(format, args...)
}

// GinLogger returns a gin middleware for logging requests
func GinLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start timer
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

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

		// Get request ID
		requestID, exists := c.Get("request_id")
		if !exists {
			requestID = "-"
		}

		// Log request
		logger.WithFields(logrus.Fields{
			"status_code":  statusCode,
			"latency":      latency,
			"client_ip":    clientIP,
			"method":       c.Request.Method,
			"path":         path,
			"query":        raw,
			"user_agent":   c.Request.UserAgent(),
			"error":        errorMessage,
			"request_id":   requestID,
			"content_type": c.ContentType(),
			"content_len":  c.Request.ContentLength,
		}).Info("Request processed")
	}
}

// GetCallerInfo returns the file name and line number of the caller
func GetCallerInfo(skip int) string {
	_, file, line, ok := runtime.Caller(skip)
	if !ok {
		return "unknown:0"
	}
	return fmt.Sprintf("%s:%d", file, line)
}

