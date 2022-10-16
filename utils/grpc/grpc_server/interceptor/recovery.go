package interceptor

import (
	"context"
	"runtime/debug"

	"go.uber.org/zap"

	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/weblazy/easy/utils/etrace"
	"github.com/weblazy/easy/utils/glog"
	"google.golang.org/grpc"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

//

func GrpcRecoveryHandler(ctx context.Context, p interface{}) (err error) {
	glog.ErrorCtx(ctx, "panic", glog.FieldTrace(etrace.ExtractTraceID(ctx)), zap.Any("err", p), zap.String("stack", string(debug.Stack())))
	// 返回一个 grpc status 错误, 像 grpc_recovery 中间件默认行为那样
	return status.Errorf(codes.Internal, "panic: %v", p)
}

func UnaryRecoveryInterceptor() grpc.UnaryServerInterceptor {
	return grpc_recovery.UnaryServerInterceptor(
		grpc_recovery.WithRecoveryHandlerContext(GrpcRecoveryHandler))
}
