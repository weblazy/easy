package ehttp

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (

	// ClientHandleCounter ...
	ClientHandleCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "",
			Name:      "http_client_handle_total",
		}, []string{"name", "method", "path", "peer", "code"})

	// ClientHandleHistogram ...
	ClientHandleHistogram = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "",
		Name:      "http_client_handle_seconds",
	}, []string{"name", "method", "path", "peer"})
)

func init() {
	prometheus.MustRegister(ClientHandleCounter)
	prometheus.MustRegister(ClientHandleHistogram)
}
