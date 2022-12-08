package interceptor

import (
	"context"

	"github.com/weblazy/easy/transport"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func GrpcHeaderCarrierInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		if md, ok := metadata.FromIncomingContext(ctx); ok {
			ctx = transport.CustomKeysMapPropagator.Extract(ctx, transport.GrpcHeaderCarrier(md))
		}
		return handler(ctx, req)
	}
}
