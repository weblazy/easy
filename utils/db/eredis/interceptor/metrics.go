package interceptor

import "github.com/prometheus/client_golang/prometheus"

var (

	// ClientHandleCounter ...
	RedisHandleCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "",
			Name:      "redis_handle_total",
		}, []string{"name", "method", "peer", "code"})

	// ClientHandleHistogram ...
	RedisHandleHistogram = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "",
		Name:      "redis_handle_seconds",
	}, []string{"name", "method", "peer"})
)

func init() {
	prometheus.MustRegister(RedisHandleCounter)
	prometheus.MustRegister(RedisHandleHistogram)
}
