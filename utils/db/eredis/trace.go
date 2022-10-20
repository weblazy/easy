package eredis

import (
	"context"

	"github.com/go-redis/redis/extra/rediscmd/v8"
	"github.com/go-redis/redis/v8"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
	"go.opentelemetry.io/otel/trace"
)

// copy from https://github.com/go-redis/redis/blob/v8/extra/redisotel/redisotel.go
// 主要做了以下修改:
//   1. 不会将命令明文当做 attr

const (
	defaultTracerName = "github.com/go-redis/redis/extra/redisotel"
)

type TracingHook struct {
	tracer trace.Tracer
	attrs  []attribute.KeyValue
}

func NewTracingHook(opts ...TraceOption) *TracingHook {
	cfg := &config{
		tp: otel.GetTracerProvider(),
		attrs: []attribute.KeyValue{
			semconv.DBSystemRedis,
		},
	}
	for _, opt := range opts {
		opt.apply(cfg)
	}

	tracer := cfg.tp.Tracer(
		defaultTracerName,
		trace.WithInstrumentationVersion("semver:"+redis.Version()),
	)
	return &TracingHook{tracer: tracer, attrs: cfg.attrs}
}

func (th *TracingHook) BeforeProcess(ctx context.Context, cmd redis.Cmder) (context.Context, error) {
	if !trace.SpanFromContext(ctx).IsRecording() {
		return ctx, nil
	}

	opts := []trace.SpanStartOption{
		trace.WithSpanKind(trace.SpanKindClient),
		trace.WithAttributes(th.attrs...),
		trace.WithAttributes(
			semconv.DBStatementKey.String(cmd.Name()),
		),
	}

	ctx, _ = th.tracer.Start(ctx, cmd.FullName(), opts...)

	return ctx, nil
}

func (th *TracingHook) AfterProcess(ctx context.Context, cmd redis.Cmder) error {
	span := trace.SpanFromContext(ctx)
	if err := cmd.Err(); err != nil {
		recordError(ctx, span, err)
	}
	span.End()
	return nil
}

func (th *TracingHook) BeforeProcessPipeline(ctx context.Context, cmds []redis.Cmder) (context.Context, error) {
	if !trace.SpanFromContext(ctx).IsRecording() {
		return ctx, nil
	}

	summary, _ := rediscmd.CmdsString(cmds)
	cmdStrs := make([]string, 0, len(cmds))
	for _, cmd := range cmds {
		cmdStrs = append(cmdStrs, cmd.Name())
	}

	opts := []trace.SpanStartOption{
		trace.WithSpanKind(trace.SpanKindClient),
		trace.WithAttributes(th.attrs...),
		trace.WithAttributes(
			semconv.DBStatementKey.StringSlice(cmdStrs),
			attribute.Int("db.redis.num_cmd", len(cmds)),
		),
	}

	ctx, _ = th.tracer.Start(ctx, "pipeline "+summary, opts...)

	return ctx, nil
}

func (th *TracingHook) AfterProcessPipeline(ctx context.Context, cmds []redis.Cmder) error {
	span := trace.SpanFromContext(ctx)
	if err := cmds[0].Err(); err != nil {
		recordError(ctx, span, err)
	}
	span.End()
	return nil
}

func recordError(ctx context.Context, span trace.Span, err error) {
	if err != redis.Nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
	}
}

type config struct {
	tp    trace.TracerProvider
	attrs []attribute.KeyValue
}

// TraceOption specifies instrumentation configuration options.
type TraceOption interface {
	apply(*config)
}

type optionFunc func(*config)

func (o optionFunc) apply(c *config) {
	o(c)
}

// WithTracerProvider specifies a tracer provider to use for creating a tracer.
// If none is specified, the global provider is used.
func WithTracerProvider(provider trace.TracerProvider) TraceOption {
	return optionFunc(func(cfg *config) {
		if provider != nil {
			cfg.tp = provider
		}
	})
}

// WithAttributes specifies additional attributes to be added to the span.
func WithAttributes(attrs ...attribute.KeyValue) TraceOption {
	return optionFunc(func(cfg *config) {
		cfg.attrs = append(cfg.attrs, attrs...)
	})
}
