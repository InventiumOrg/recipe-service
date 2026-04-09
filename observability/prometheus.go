package observability

import (
	"log/slog"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// PrometheusMetrics holds all Prometheus metrics for the recipe service
type PrometheusMetrics struct {
	// HTTP metrics
	HTTPRequestsTotal       *prometheus.CounterVec
	HTTPRequestDuration     *prometheus.HistogramVec
	HTTPRequestsInFlight    prometheus.Gauge
	HTTPResponseStatusTotal *prometheus.CounterVec

	// Database metrics
	DBConnectionsActive prometheus.Gauge
	DBOperationDuration *prometheus.HistogramVec
	DBOperationErrors   *prometheus.CounterVec

	// Business metrics
	RecipeOperationsTotal  *prometheus.CounterVec
	RecipesActive          prometheus.Gauge
	AuthenticationAttempts *prometheus.CounterVec
}

// NewPrometheusMetrics creates and registers all Prometheus metrics
func NewPrometheusMetrics(serviceName string) *PrometheusMetrics {
	metrics := &PrometheusMetrics{
		// HTTP metrics following Prometheus naming conventions
		HTTPRequestsTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "http_requests_total",
				Help: "Total number of HTTP requests",
			},
			[]string{"method", "endpoint", "status_code"},
		),
		HTTPRequestDuration: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "http_request_duration_seconds",
				Help:    "HTTP request duration in seconds",
				Buckets: prometheus.DefBuckets,
			},
			[]string{"method", "endpoint"},
		),
		HTTPRequestsInFlight: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Name: "http_requests_in_flight",
				Help: "Current number of HTTP requests being processed",
			},
		),
		HTTPResponseStatusTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "http_response_status_total",
				Help: "Total number of HTTP responses by status class",
			},
			[]string{"method", "endpoint", "status_class"},
		),

		// Database metrics
		DBConnectionsActive: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Name: "database_connections_active",
				Help: "Number of active database connections",
			},
		),
		DBOperationDuration: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "database_operation_duration_seconds",
				Help:    "Database operation duration in seconds",
				Buckets: []float64{.001, .005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5},
			},
			[]string{"operation", "table"},
		),
		DBOperationErrors: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "database_operation_errors_total",
				Help: "Total number of database operation errors",
			},
			[]string{"operation", "table", "error_type"},
		),

		// Business metrics specific to recipe service
		RecipeOperationsTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "recipe_operations_total",
				Help: "Total number of recipe operations",
			},
			[]string{"operation", "recipe_name"},
		),
		RecipesActive: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Name: "recipes_active",
				Help: "Current number of active recipes (best-effort)",
			},
		),
		AuthenticationAttempts: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "authentication_attempts_total",
				Help: "Total number of authentication attempts",
			},
			[]string{"status", "method"},
		),
	}

	// Register all metrics with Prometheus
	prometheus.MustRegister(
		metrics.HTTPRequestsTotal,
		metrics.HTTPRequestDuration,
		metrics.HTTPRequestsInFlight,
		metrics.HTTPResponseStatusTotal,
		metrics.DBConnectionsActive,
		metrics.DBOperationDuration,
		metrics.DBOperationErrors,
		metrics.RecipeOperationsTotal,
		metrics.RecipesActive,
		metrics.AuthenticationAttempts,
	)

	slog.Info("Prometheus metrics registered", slog.String("service", serviceName))
	return metrics
}

// getStatusClass converts HTTP status code to status class (2xx, 4xx, 5xx, etc.)
func getStatusClass(statusCode int) string {
	switch {
	case statusCode >= 200 && statusCode < 300:
		return "2xx"
	case statusCode >= 300 && statusCode < 400:
		return "3xx"
	case statusCode >= 400 && statusCode < 500:
		return "4xx"
	case statusCode >= 500:
		return "5xx"
	default:
		return "1xx"
	}
}

// PrometheusMiddleware creates a Gin middleware for collecting HTTP metrics
func (m *PrometheusMetrics) PrometheusMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Skip metrics collection for the metrics endpoint itself
		if c.Request.URL.Path == "/metrics" {
			c.Next()
			return
		}

		start := time.Now()

		// Increment in-flight requests
		m.HTTPRequestsInFlight.Inc()
		defer m.HTTPRequestsInFlight.Dec()

		// Process request
		c.Next()

		// Calculate duration
		duration := time.Since(start).Seconds()

		// Get route pattern
		route := c.FullPath()
		if route == "" {
			route = "unknown"
		}

		// Record metrics
		statusCode := c.Writer.Status()
		statusClass := getStatusClass(statusCode)

		m.HTTPRequestsTotal.WithLabelValues(
			c.Request.Method,
			route,
			strconv.Itoa(statusCode),
		).Inc()

		m.HTTPRequestDuration.WithLabelValues(
			c.Request.Method,
			route,
		).Observe(duration)

		m.HTTPResponseStatusTotal.WithLabelValues(
			c.Request.Method,
			route,
			statusClass,
		).Inc()
	}
}

// RecordDBOperation records database operation metrics
func (m *PrometheusMetrics) RecordDBOperation(operation, table string, duration time.Duration, err error) {
	m.DBOperationDuration.WithLabelValues(operation, table).Observe(duration.Seconds())

	if err != nil {
		errorType := "unknown"
		m.DBOperationErrors.WithLabelValues(operation, table, errorType).Inc()
	}
}

// RecordRecipeOperation records business-specific recipe operations.
func (m *PrometheusMetrics) RecordRecipeOperation(operation, recipeName string) {
	if recipeName == "" {
		recipeName = "unknown"
	}
	m.RecipeOperationsTotal.WithLabelValues(operation, recipeName).Inc()
}

// UpdateRecipesCount updates the current count of active recipes (best-effort).
func (m *PrometheusMetrics) UpdateRecipesCount(delta float64) {
	m.RecipesActive.Add(delta)
}

// RecordAuthAttempt records authentication attempts
func (m *PrometheusMetrics) RecordAuthAttempt(status, method string) {
	m.AuthenticationAttempts.WithLabelValues(status, method).Inc()
}

// UpdateDBConnections updates the database connections gauge
func (m *PrometheusMetrics) UpdateDBConnections(count float64) {
	m.DBConnectionsActive.Set(count)
}

// SetupPrometheusEndpoint adds the /metrics endpoint to the Gin router
func SetupPrometheusEndpoint(router *gin.Engine) {
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))
	slog.Info("Prometheus metrics endpoint configured at /metrics")
}

// WithDBMetrics wraps a database operation with automatic metrics collection
func (m *PrometheusMetrics) WithDBMetrics(operation, table string, fn func() error) error {
	start := time.Now()
	err := fn()
	m.RecordDBOperation(operation, table, time.Since(start), err)
	return err
}

// WithRecipeMetrics wraps a recipe operation with automatic metrics collection
func (m *PrometheusMetrics) WithRecipeMetrics(operation, recipeName string, fn func() error) error {
	err := fn()
	if err == nil {
		m.RecordRecipeOperation(operation, recipeName)
	}
	return err
}
