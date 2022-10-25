package interceptor

import (
	"context"
	"time"

	"go.uber.org/zap"

	"github.com/google/uuid"
	"github.com/weblazy/easy/utils/elog"
	"github.com/weblazy/easy/utils/elog/ezap"
	"github.com/weblazy/easy/utils/etrace"
	"github.com/weblazy/easy/utils/grpc/grpc_server/grpc_server_config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func init() {
	logger := ezap.NewFileEzap(grpc_server_config.PkgName)
	elog.SetLogger(grpc_server_config.PkgName, logger)
}

func GrpcLogger(logConf *elog.LogConf) grpc.UnaryServerInterceptor {
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
		elog.SetContextLog(ctx, logConf)

		reqLabel, mdLabel := zap.Any("req", req), zap.Any("metadata", md)

		resp, err = handler(ctx, req)
		ctx = elog.SetLogerName(ctx, grpc_server_config.PkgName)
		if err != nil {
			elog.ErrorCtx(ctx, "grpc_server", reqLabel, mdLabel, elog.FieldError(err), elog.FieldCost(time.Since(start)))
		} else {
			elog.InfoCtx(ctx, "grpc_server", reqLabel, mdLabel, zap.Any("res", resp), elog.FieldCost(time.Since(start)))
		}

		return resp, err
	}
}

func GrpcLoggerLite(logConf elog.LogConf) grpc.UnaryServerInterceptor {
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
