package interceptor

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	ServerHandledCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "fgo",
			Name:      "grpc_server_handled_total",
			Help:      "Total number of RPCs completed on the server, regardless of success or failure.",
		}, []string{"grpc_type", "method", "code", "uniform_code"})

	ServerHandledHistogram = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "fgo",
			Name:      "grpc_server_handling_seconds",
			Help:      "Histogram of response latency (seconds) of gRPC that had been application-level handled by the server.",
		}, []string{"grpc_type", "method"})
)
