package interceptor

import (
	"context"
	"time"

	"google.golang.org/grpc"
)

// copy from go-zero https://github.com/zeromicro/go-zero/blob/2732d3cdae5bf35dc07e926d3b5ed35e3c506393/zrpc/internal/clientinterceptors/timeoutinterceptor.go

type contextKeyType struct{}

var ctxKey contextKeyType

// TimeoutInterceptor is an interceptor that controls timeout.
func TimeoutInterceptor(timeout time.Duration) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		// fix closure error
		t := timeout

		if v, ok := ctx.Value(ctxKey).(time.Duration); ok {
			t = v
		}

		if v, ok := getForceTimeout(opts); ok {
			t = v
		}

		if t <= 0 {
			return invoker(ctx, method, req, reply, cc, opts...)
		}

		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, t)
		defer cancel()

		return invoker(ctx, method, req, reply, cc, opts...)
	}
}

// CallOption is a grpc.CallOption that is local to timeout interceptor.
type CallOption struct {
	grpc.EmptyCallOption

	forceTimeout time.Duration
}

// WithForceTimeout sets the RPC timeout for this call only.
func WithForceTimeout(forceTimeout time.Duration) CallOption {
	return CallOption{forceTimeout: forceTimeout}
}

func getForceTimeout(callOptions []grpc.CallOption) (time.Duration, bool) {
	for _, opt := range callOptions {
		if co, ok := opt.(CallOption); ok {
			return co.forceTimeout, true
		}
	}

	return 0, false
}

// ForceTimeout force set timeout for this rpc call
// Deprecated: use WithForceTimeout
func ForceTimeout(ctx context.Context, timeout time.Duration) context.Context {
	if timeout <= 0 {
		return ctx
	}

	return context.WithValue(ctx, ctxKey, timeout)
}
