package etrace

import (
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

var (
	// HTTPMiddleware is alias for telhttp.NewHandler.
	HTTPMiddleware = otelhttp.NewHandler
)

// HTTPTransport is alias for otelhttp.NewTransport.
var HTTPTransport = otelhttp.NewTransport
