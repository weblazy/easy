package interceptor

import (
	"context"

	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/weblazy/easy/utils/eerror"
	"github.com/weblazy/easy/utils/fmetric"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

var (
	// ClientHandleCounter ...
	ClientHandleCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "grpc_client_handle_total",
		}, []string{"type", "name", "method", "peer", "code"})

	// ClientHandleHistogram ...
	ClientHandleHistogram = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "grpc_client_handle_seconds",
		}, []string{"type", "name", "method", "peer"})
)

var (
	ClientWithBizHandledCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "monitor_grpc_client_accept_result_total",
			Help: "Total number of RPCs accept on the client, regardless of success or failure.",
		}, []string{"grpc_type", "grpc_service", "grpc_method", "grpc_code", "result_code"})
)

func MetricUnaryClientInterceptor(successCodes []string) grpc.UnaryClientInterceptor {
	grpc_prometheus.EnableClientHandlingTimeHistogram()
	originMw := grpc_prometheus.UnaryClientInterceptor
	extractor := eerror.ExtractBizCode(successCodes)

	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		err := originMw(ctx, method, req, reply, cc, invoker, opts...)
		st, _ := status.FromError(err)

		bizCode, ok := extractor(reply, err)
		if ok {
			service, method := fmetric.SplitGrpcMethodName(method)
			ClientWithBizHandledCounter.WithLabelValues("unary", service, method, st.Code().String(), bizCode).Inc()
		}
		return err
	}
}
