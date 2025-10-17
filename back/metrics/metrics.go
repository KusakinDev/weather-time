package metrics

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	httpRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "app",
			Subsystem: "http",
			Name:      "requests_total",
			Help:      "Total number of HTTP requests.",
		},
		[]string{"method", "route", "status"}, // status = numeric code (e.g., "200")
	)

	httpRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "app",
			Subsystem: "http",
			Name:      "request_duration_seconds",
			Help:      "HTTP request latency.",
			// под типичный веб-профиль; скорректируй под себя
			Buckets: []float64{0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10},
		},
		[]string{"method", "route", "status"},
	)

	httpInFlight = promauto.NewGauge(
		prometheus.GaugeOpts{
			Namespace: "app",
			Subsystem: "http",
			Name:      "in_flight_requests",
			Help:      "In-flight HTTP requests.",
		},
	)
)

// GinMiddleware — вешай через r.Use(Middleware()) СНИЗУ ВСЕХ recovery/logger,
// чтобы route = c.FullPath() уже был известен
func GinMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		httpInFlight.Inc()
		defer httpInFlight.Dec()

		c.Next()

		route := c.FullPath()
		if route == "" {
			route = "unknown"
		}
		status := strconv.Itoa(c.Writer.Status())

		httpRequestsTotal.WithLabelValues(
			c.Request.Method, route, status,
		).Inc()

		lat := time.Since(start).Seconds()
		httpRequestDuration.WithLabelValues(
			c.Request.Method, route, status,
		).Observe(lat)
	}
}
