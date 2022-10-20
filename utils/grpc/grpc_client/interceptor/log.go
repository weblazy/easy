package interceptor

import (
	"context"
	"net/http"
	"time"

	"github.com/weblazy/easy/utils/ecodes"
	"github.com/weblazy/easy/utils/elog"
	"github.com/weblazy/easy/utils/etrace"
	"github.com/weblazy/easy/utils/transport"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type LogConf struct {
	EnableTraceInterceptor     bool
	EnableAccessInterceptorReq bool          // 是否开启记录请求参数，默认开启
	EnableAccessInterceptorRes bool          // 是否开启记录响应参数，默认开启
	SlowLogThreshold           time.Duration // 慢日志记录的阈值，默认600ms
	EnableAccessInterceptor    bool          // 是否开启记录请求数据，默认开启
}

// loggerUnaryClientInterceptor returns log interceptor for logging
func LoggerUnaryClientInterceptor(config *LogConf) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, res interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		beg := time.Now()
		var fields = make([]zap.Field, 0, 20)

		var md metadata.MD
		// try to append custom metadata to client request metadata
		if m, ok := metadata.FromOutgoingContext(ctx); ok {
			md = m.Copy()
		} else {
			md = metadata.MD{}
		}

		transport.CustomKeysMapPropagator.Inject(ctx, transport.GrpcHeaderCarrier(md))
		ctx = metadata.NewOutgoingContext(ctx, md)

		err := invoker(ctx, method, req, res, cc, opts...)
		duration := time.Since(beg)
		spbStatus := ecodes.Convert(err)
		httpStatusCode := ecodes.GrpcToHTTPStatusCode(spbStatus.Code())

		fields = append(fields,
			zap.String("type", "unary"),
			zap.Int64("code", int64(spbStatus.Code())),
			zap.Int64("uniformCode", int64(httpStatusCode)),
			zap.String("description", spbStatus.Message()),
			elog.FieldMethod(method),
			elog.FieldCost(duration),
			elog.FieldName(cc.Target()),
		)

		span := trace.SpanFromContext(ctx)
		// add custom metadata to trace fields
		for k, v := range transport.GetMapFromContext(ctx) {
			span.SetAttributes(attribute.String(k, v))
		}

		// 开启了链路，那么就记录链路id
		if config.EnableTraceInterceptor && etrace.IsGlobalTracerRegistered() {
			fields = append(fields, elog.FieldTrace(etrace.ExtractTraceID(ctx)))
		}

		if config.EnableAccessInterceptorReq {
			fields = append(fields, zap.Any("req", req))
		}
		if config.EnableAccessInterceptorRes {
			fields = append(fields, zap.Any("res", res))
		}

		if config.SlowLogThreshold > time.Duration(0) && duration > config.SlowLogThreshold {
			elog.WarnCtx(ctx, "slow", fields...)
		}

		if err != nil {
			fields = append(fields, elog.FieldEvent("error"), elog.FieldError(err))
			// 只记录系统级别错误
			if httpStatusCode >= http.StatusInternalServerError {
				// 只记录系统级别错误
				elog.ErrorCtx(ctx, "access", fields...)
				return err
			}
			// 业务报错只做warning
			elog.WarnCtx(ctx, "access", fields...)
			return err
		}

		if config.EnableAccessInterceptor {
			fields = append(fields, elog.FieldEvent("normal"))
			elog.InfoCtx(ctx, "access", fields...)
		}
		return nil
	}
}
