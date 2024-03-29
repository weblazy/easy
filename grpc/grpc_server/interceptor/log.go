package interceptor

import (
	"context"
	"sync"
	"time"

	"go.uber.org/zap"

	"github.com/google/uuid"
	"github.com/weblazy/easy/elog"
	"github.com/weblazy/easy/etrace"
	"github.com/weblazy/easy/grpc/grpc_server/grpc_server_config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

var once sync.Once

func GrpcLogger(config *grpc_server_config.Config) grpc.UnaryServerInterceptor {
	once.Do(config.InitLogger)
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
		fields := make([]zap.Field, 0)
		fields = append(fields, elog.FieldMethod(info.FullMethod), elog.FieldReq(req), zap.Any("metadata", md))

		resp, err = handler(ctx, req)
		ctx = elog.SetLogerName(ctx, grpc_server_config.PkgName)
		if err != nil {
			fields = append(fields, elog.FieldError(err), elog.FieldDuration(time.Since(start)))
			elog.ErrorCtx(ctx, grpc_server_config.PkgName, fields...)
		} else {
			fields = append(fields, elog.FieldResp(resp), elog.FieldDuration(time.Since(start)))
			elog.InfoCtx(ctx, grpc_server_config.PkgName, fields...)
		}
		return resp, err
	}
}

// func GrpcLoggerLite() grpc.UnaryServerInterceptor {
// 	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
// 		// otel trace
// 		traceId := etrace.ExtractTraceID(ctx)

// 		md, _ := metadata.FromIncomingContext(ctx)
// 		// 尝试获取网关 traceid
// 		if traceId == "" {
// 			v := md.Get("traceid")
// 			if len(v) > 0 {
// 				traceId = v[0]
// 			}
// 		}

// 		// 服务内部生成
// 		if traceId == "" {
// 			traceId = uuid.NewString()
// 		}
// 		logConf.Name = "server.grpc"
// 		logConf.Fields = append(logConf.Fields, zap.String("trace_id", traceId), zap.String("method", info.FullMethod))
// 		return handler(ctx, req)
// 	}
// }
