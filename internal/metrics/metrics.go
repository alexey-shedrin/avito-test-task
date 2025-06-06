package metrics

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var requestCount = promauto.NewCounterVec(
	prometheus.CounterOpts{
		Name: "http.requests.total",
		Help: "Total number of HTTP requests by handler, method and status code",
	},
	[]string{"handler", "method", "code"},
)

var responseTime = promauto.NewHistogramVec(
	prometheus.HistogramOpts{
		Name:    "http.response.time.seconds",
		Help:    "Histogram of response times for HTTP requests",
		Buckets: prometheus.DefBuckets,
	},
	[]string{"handler", "method"},
)

var createPvzCount = promauto.NewCounter(
	prometheus.CounterOpts{
		Name: "created.pvz.total",
		Help: "Total number of created PVZ",
	},
)

func CreatePVZ() {
	createPvzCount.Inc()
}

var createdReceptionCount = promauto.NewCounter(
	prometheus.CounterOpts{
		Name: "created.reception.total",
		Help: "Total number of created Receptions",
	},
)

func CreateReception() {
	createdReceptionCount.Inc()
}

var addedProductCount = promauto.NewCounter(
	prometheus.CounterOpts{
		Name: "added.product.total",
		Help: "Total number of added Products",
	},
)

func AddProduct() {
	addedProductCount.Inc()
}

func DeleteProduct() {
	addedProductCount.Desc()
}

func StartMetricsServer(port string) {
	http.Handle("/metrics", promhttp.Handler())
	if err := http.ListenAndServe(fmt.Sprintf(":%s", port), nil); err != nil {
		log.Fatalf("metrics server failed: %v", err)
	}
}
func GetMetricsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		duration := time.Since(start).Seconds()

		status := c.Writer.Status()
		path := sanitizePath(c.Request.URL.Path)

		requestCount.WithLabelValues(path, c.Request.Method, getStatusCode(status)).Inc()
		responseTime.WithLabelValues(path, c.Request.Method).Observe(duration)
	}
}

func sanitizePath(p string) string {
	if strings.HasPrefix(p, "/pvz/") {
		return "/pvz/:id"
	}
	if strings.HasPrefix(p, "/reception/") {
		return "/reception/:id"
	}
	return p
}

func getStatusCode(code int) string {
	return strconv.Itoa(code)
}
