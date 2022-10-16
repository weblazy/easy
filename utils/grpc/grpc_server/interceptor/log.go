package interceptor

import (
	"context"
	"time"

	"go.uber.org/zap"

	"github.com/google/uuid"
	"github.com/weblazy/easy/utils/etrace"
	"github.com/weblazy/easy/utils/glog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func GrpcLogger(logConf *glog.LogConf) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		start := time.Now()

		// otel trace
		traceId := etrace.ExtractTraceID(ctx)

		md, _ := metadata.FromIncomingContext(ctx)
		// 尝试获取网关 traceid
		if traceId == "" {
			v := md.Get("traceid")
			if len(v) > 0 {
				traceId = v[0]
			}
		}

		// 服务内部生成
		if traceId == "" {
			traceId = uuid.NewString()
		}
		logConf.Name = "server.grpc"
		logConf.Labels = append(logConf.Labels, zap.String("trace_id", traceId), zap.String("method", info.FullMethod))

		// // set new logger to context
		// ctx = blog.NewContext(ctx, logger)
		glog.SetContextLog(ctx, logConf)

		reqLabel, mdLabel := zap.Any("request", req), zap.Any("metadata", md)

		resp, err = handler(ctx, req)

		if err != nil {
			glog.ErrorCtx(ctx, "grpc log", reqLabel, mdLabel, glog.FieldError(err), glog.FieldCost(time.Since(start)))
		} else {
			glog.InfoCtx(ctx, "grpc log", reqLabel, mdLabel, zap.Any("response", resp), glog.FieldCost(time.Since(start)))
		}

		return resp, err
	}
}

func GrpcLoggerLite(logConf glog.LogConf) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		// otel trace
		traceId := etrace.ExtractTraceID(ctx)

		md, _ := metadata.FromIncomingContext(ctx)
		// 尝试获取网关 traceid
		if traceId == "" {
			v := md.Get("traceid")
			if len(v) > 0 {
				traceId = v[0]
			}
		}

		// 服务内部生成
		if traceId == "" {
			traceId = uuid.NewString()
		}
		logConf.Name = "server.grpc"
		logConf.Labels = append(logConf.Labels, zap.String("trace_id", traceId), zap.String("method", info.FullMethod))

		return handler(ctx, req)
	}
}
