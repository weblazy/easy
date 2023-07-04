package interceptor

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/weblazy/easy/http/http_server/http_server_config"
)

var (

	// ServerHandleCounter ...
	ServerHandleCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "",
			Name:      "http_server_handle_total",
		}, []string{"name", "method", "path", "host", "code"})

	// ServerHandleHistogram ...
	ServerHandleHistogram = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "",
		Name:      "http_server_handle_seconds",
	}, []string{"name", "method", "path", "host"})
)

func init() {
	prometheus.MustRegister(ServerHandleCounter)
	prometheus.MustRegister(ServerHandleHistogram)
}
func MetricInterceptor(cfg *http_server_config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		if cfg.MetricPathRewriter == nil {
			cfg.MetricPathRewriter = http_server_config.DefaultMetricPathRewriter
		}
		ServerHandleCounter.WithLabelValues(cfg.Name, c.Request.Method, cfg.MetricPathRewriter(c.Request.URL.Path), c.Request.URL.Host, strconv.Itoa(c.Writer.Status())).Inc()
		ServerHandleHistogram.WithLabelValues(cfg.Name, c.Request.Method, cfg.MetricPathRewriter(c.Request.URL.Path), c.Request.URL.Host).Observe(time.Since(GetStartTime(c.Request.Context())).Seconds())
	}
}
