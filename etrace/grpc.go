package etrace

import "go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"

var (
	// UnaryServerInterceptor is alias for otelgrpc.UnaryServerInterceptor.
	UnaryServerInterceptor = otelgrpc.UnaryServerInterceptor
	// StreamServerInterceptor is alias for otelgrpc.StreamServerInterceptor.
	StreamServerInterceptor = otelgrpc.StreamServerInterceptor
)

var (
	// UnaryClientInterceptor is alias for  otelgrpc.UnaryClientInterceptor.
	UnaryClientInterceptor = otelgrpc.UnaryClientInterceptor
	// StreamClientInterceptor is alias for otelgrpc.StreamClientInterceptor.
	StreamClientInterceptor = otelgrpc.StreamClientInterceptor
)
