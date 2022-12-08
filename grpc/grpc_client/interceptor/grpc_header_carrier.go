package interceptor

import (
	"context"

	"github.com/weblazy/easy/transport"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// GrpcHeaderCarrierInterceptor
func GrpcHeaderCarrierInterceptor() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		var md metadata.MD
		// try to append custom metadata to client request metadata
		if m, ok := metadata.FromOutgoingContext(ctx); ok {
			md = m.Copy()
		} else {
			md = metadata.MD{}
		}
		transport.CustomKeysMapPropagator.Inject(ctx, transport.GrpcHeaderCarrier(md))
		ctx = metadata.NewOutgoingContext(ctx, md)
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}
