package interceptor

import (
	"github.com/go-resty/resty/v2"
	"github.com/spf13/cast"
	"github.com/weblazy/easy/http/http_client/http_client_config"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// Deprecated: use otel http transport
func TraceInterceptor(name string, cfg *http_client_config.Config) (resty.RequestMiddleware, resty.ResponseMiddleware, resty.ErrorHook) { //nolint
	tracer := otel.Tracer("")

	beforeFn := func(cli *resty.Client, req *resty.Request) error {
		ctx, span := tracer.Start(req.Context(), req.Method, trace.WithSpanKind(trace.SpanKindClient))

		span.SetAttributes(
			attribute.String("peer.service", name),
			attribute.String("http.method", req.Method),
			attribute.String("http.url", req.URL),
		)

		req.SetContext(ctx)
		return nil
	}

	afterFn := func(cli *resty.Client, res *resty.Response) error {
		span := trace.SpanFromContext(res.Request.Context())
		span.SetAttributes(
			attribute.String("http.status_code", cast.ToString(res.StatusCode())),
		)

		span.End()
		return nil
	}

	errorFn := func(req *resty.Request, err error) {
		span := trace.SpanFromContext(req.Context())

		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
		}

		span.End()
	}
	return beforeFn, afterFn, errorFn
}
