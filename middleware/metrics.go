package middleware

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// httpRequestsTotal HTTP请求总数
	httpRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "path", "status"},
	)

	// httpRequestDuration HTTP请求延迟
	httpRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "HTTP request duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path"},
	)

	// activeRequests 当前活跃请求数
	activeRequests = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "http_active_requests",
			Help: "Number of active HTTP requests",
		},
	)

	// productOperationsTotal 商品操作总数
	productOperationsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "product_operations_total",
			Help: "Total number of product operations",
		},
		[]string{"operation", "status"},
	)

	// orderOperationsTotal 订单操作总数
	orderOperationsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "order_operations_total",
			Help: "Total number of order operations",
		},
		[]string{"operation", "status"},
	)

	// cacheHitsTotal 缓存命中总数
	cacheHitsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "cache_hits_total",
			Help: "Total number of cache hits",
		},
		[]string{"cache_type", "result"},
	)

	// dbQueryDuration 数据库查询延迟
	dbQueryDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "db_query_duration_seconds",
			Help:    "Database query duration in seconds",
			Buckets: []float64{.001, .005, .01, .025, .05, .1, .25, .5, 1},
		},
		[]string{"operation", "table"},
	)
)

// MetricsMiddleware Prometheus指标中间件
func MetricsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.FullPath()
		if path == "" {
			path = "unknown"
		}

		// 增加活跃请求计数
		activeRequests.Inc()

		// 处理请求
		c.Next()

		// 计算延迟
		duration := time.Since(start).Seconds()
		status := strconv.Itoa(c.Writer.Status())

		// 记录指标
		httpRequestsTotal.WithLabelValues(c.Request.Method, path, status).Inc()
		httpRequestDuration.WithLabelValues(c.Request.Method, path).Observe(duration)

		// 减少活跃请求计数
		activeRequests.Dec()
	}
}

// RecordProductOperation 记录商品操作
func RecordProductOperation(operation string, success bool) {
	status := "success"
	if !success {
		status = "failure"
	}
	productOperationsTotal.WithLabelValues(operation, status).Inc()
}

// RecordOrderOperation 记录订单操作
func RecordOrderOperation(operation string, success bool) {
	status := "success"
	if !success {
		status = "failure"
	}
	orderOperationsTotal.WithLabelValues(operation, status).Inc()
}

// RecordCacheHit 记录缓存命中
func RecordCacheHit(cacheType string, hit bool) {
	result := "hit"
	if !hit {
		result = "miss"
	}
	cacheHitsTotal.WithLabelValues(cacheType, result).Inc()
}

// RecordDBQuery 记录数据库查询
func RecordDBQuery(operation, table string, duration time.Duration) {
	dbQueryDuration.WithLabelValues(operation, table).Observe(duration.Seconds())
}
