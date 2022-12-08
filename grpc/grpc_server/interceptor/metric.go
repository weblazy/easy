package interceptor

import (
	"context"

	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/weblazy/easy/eerror"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"

	"github.com/weblazy/easy/fmetric"
)

// 目前图表只用到了 resultCode 和 app 字段(收集时自动注入)

var (
	ServerWithBizHandledCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: fmetric.DefaultNamespace,
			Name:      "monitor_grpc_server_result_total",
			Help:      "Total number of RPCs completed on the server, regardless of success or failure.",
		}, []string{"grpc_type", "grpc_service", "grpc_method", "grpc_code", "result_code"})
)

func MetricUnaryServerInterceptor(successCodes []string) grpc.UnaryServerInterceptor {
	grpc_prometheus.EnableHandlingTimeHistogram()
	originMw := grpc_prometheus.UnaryServerInterceptor
	extractor := eerror.ExtractBizCode(successCodes)

	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		resp, err = originMw(ctx, req, info, handler)
		st, _ := status.FromError(err)

		bizCode, ok := extractor(resp, err)
		if ok {
			service, method := fmetric.SplitGrpcMethodName(info.FullMethod)
			ServerWithBizHandledCounter.WithLabelValues("unary", service, method, st.Code().String(), bizCode).Inc()
		}
		return
	}
}
