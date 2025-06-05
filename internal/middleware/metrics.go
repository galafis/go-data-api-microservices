package middleware

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// httpRequestsTotal counts the number of HTTP requests
	httpRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "path", "status"},
	)

	// httpRequestDuration tracks the duration of HTTP requests
	httpRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "HTTP request duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path", "status"},
	)

	// httpRequestSize tracks the size of HTTP requests
	httpRequestSize = promauto.NewSummaryVec(
		prometheus.SummaryOpts{
			Name:       "http_request_size_bytes",
			Help:       "HTTP request size in bytes",
			Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
		},
		[]string{"method", "path"},
	)

	// httpResponseSize tracks the size of HTTP responses
	httpResponseSize = promauto.NewSummaryVec(
		prometheus.SummaryOpts{
			Name:       "http_response_size_bytes",
			Help:       "HTTP response size in bytes",
			Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
		},
		[]string{"method", "path", "status"},
	)

	// activeRequests tracks the number of active HTTP requests
	activeRequests = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "http_active_requests",
			Help: "Number of active HTTP requests",
		},
	)
)

// Metrics is a middleware that collects metrics for HTTP requests
func Metrics() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Increment active requests
		activeRequests.Inc()
		defer activeRequests.Dec()

		// Start timer
		start := time.Now()

		// Get request size
		requestSize := float64(c.Request.ContentLength)
		if requestSize < 0 {
			requestSize = 0
		}

		// Track request size
		httpRequestSize.WithLabelValues(
			c.Request.Method,
			c.FullPath(),
		).Observe(requestSize)

		// Process request
		c.Next()

		// Calculate duration
		duration := time.Since(start).Seconds()

		// Get status code
		status := strconv.Itoa(c.Writer.Status())

		// Track metrics
		httpRequestsTotal.WithLabelValues(
			c.Request.Method,
			c.FullPath(),
			status,
		).Inc()

		httpRequestDuration.WithLabelValues(
			c.Request.Method,
			c.FullPath(),
			status,
		).Observe(duration)

		httpResponseSize.WithLabelValues(
			c.Request.Method,
			c.FullPath(),
			status,
		).Observe(float64(c.Writer.Size()))
	}
}

