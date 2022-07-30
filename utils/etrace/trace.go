package etrace

import (
	"context"
	"strings"

	"go.opentelemetry.io/contrib/propagators/jaeger"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

const (
	fixSpanIDPrefix = "0000000000000000"
)

type registeredTracer struct {
	isRegistered bool
}

var (
	globalTracer = registeredTracer{false}
)

func SetGlobalTracer(tp trace.TracerProvider) {
	globalTracer = registeredTracer{true}
	otel.SetTracerProvider(tp)
	// use jaeger propagator, header uber-trace-id
	otel.SetTextMapPropagator(jaeger.Jaeger{})
}

// IsGlobalTracerRegistered returns a `bool` to indicate if a tracer has been globally registered.
func IsGlobalTracerRegistered() bool {
	return globalTracer.isRegistered
}

// ExtractTraceID HTTP使用request.Context，不要使用错了.
func ExtractTraceID(ctx context.Context) string {
	if !IsGlobalTracerRegistered() {
		return ""
	}
	span := trace.SpanContextFromContext(ctx)
	if span.HasTraceID() {
		sp := span.TraceID().String()
		// https://github.com/open-telemetry/opentelemetry-go/issues/686
		// remove left padding for 64-bit TraceIDs
		return strings.TrimPrefix(sp, fixSpanIDPrefix)
	}
	return ""
}
