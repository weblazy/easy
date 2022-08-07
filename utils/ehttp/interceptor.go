package ehttp

import (
	"context"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/spf13/cast"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"

	"github.com/weblazy/easy/utils/etrace"
	"github.com/weblazy/easy/utils/glog"
	"go.uber.org/zap"
)

// type interceptor func(name string, cfg *Config, logger *blog.Logger) (resty.RequestMiddleware, resty.ResponseMiddleware, resty.ErrorHook)

type MetricPathRewriter func(origin string) string

func NoopMetricPathRewriter(origin string) string {
	return origin
}

func AddInterceptors(client *resty.Client, onBefore resty.RequestMiddleware, onAfter resty.ResponseMiddleware, onErr resty.ErrorHook) {
	if onBefore != nil {
		client.OnBeforeRequest(onBefore)
	}
	if onAfter != nil {
		client.OnAfterResponse(onAfter)
	}
	if onErr != nil {
		client.OnError(onErr)
	}
}

func logAccess(name string, config *Config, req *resty.Request, res *resty.Response, err error) {
	rr := req.RawRequest
	var url1, host string
	// 修复err 不是 *resty.ResponseError错误的时候，可能为nil
	if rr != nil {
		url1 = rr.URL.RequestURI()
		host = rr.URL.Host
	} else { // RawRequest 不一定总是有
		u, err2 := url.Parse(req.URL)
		if err2 == nil {
			url1 = u.RequestURI()
			host = u.Host
		}
	}

	fullMethod := req.Method + "." + url1 // GET./hello
	var cost = time.Since(beg(req.Context()))
	var respBody string
	if res != nil {
		respBody = string(res.Body())
	}

	var fields = make([]zap.Field, 0, 15)
	fields = append(fields,
		zap.String("method", fullMethod),
		zap.String("name", name),
		zap.Float64("cost", float64(cost.Microseconds())/1000),
		zap.String("host", host),
	)

	// 开启了链路，那么就记录链路id
	// todo
	if config.EnableTraceInterceptor && etrace.IsGlobalTracerRegistered() {
		fields = append(fields, zap.String("trace_id", etrace.ExtractTraceID(req.Context())))
	}

	if config.EnableAccessInterceptorReq {
		reqMap := map[string]interface{}{
			"url":  req.URL,
			"body": req.Body,
		}

		if config.EnableAccessInterceptorReqHeader {
			reqMap["header"] = req.Header
		}

		fields = append(fields, zap.Any("req", reqMap))
	}

	if config.EnableAccessInterceptorRes {
		resMap := make(map[string]interface{}, 3)
		resMap["body"] = respBody

		// 处理 res 为空时空指针错误
		if res != nil {
			resMap["header"] = res.Header()
			resMap["status_code"] = res.StatusCode()
		}

		fields = append(fields, zap.Any("res", resMap))
	}

	if config.SlowLogThreshold > time.Duration(0) && cost > config.SlowLogThreshold {
		glog.InfoCtx(req.Context(), "slow", fields...)
	}

	if err != nil {
		fields = append(fields, zap.String("event", "error"), zap.Error(err))
		if res == nil {
			// 无 res 的是连接超时等系统级错误
			glog.ErrorCtx(req.Context(), "access", fields...)
			return
		}
		glog.InfoCtx(req.Context(), "access", fields...)
		return
	}

	if config.EnableAccessInterceptor {
		fields = append(fields, zap.String("event", "normal"))
		glog.InfoCtx(req.Context(), "access", fields...)
	}
}

// https://stackoverflow.com/questions/40891345/fix-should-not-use-basic-type-string-as-key-in-context-withvalue-golint
// https://blog.golang.org/context#TOC_3.2.
// https://golang.org/pkg/context/#WithValue ，这边文章说明了用struct，可以避免分配
type begKey struct{}

func beg(ctx context.Context) time.Time {
	begTime, _ := ctx.Value(begKey{}).(time.Time)
	return begTime
}

func fixedInterceptor(name string, config *Config) (resty.RequestMiddleware, resty.ResponseMiddleware, resty.ErrorHook) {
	return func(cli *resty.Client, req *resty.Request) error {
		req.SetContext(context.WithValue(req.Context(), begKey{}, time.Now()))
		return nil
	}, nil, nil
}

func logInterceptor(name string, config *Config) (resty.RequestMiddleware, resty.ResponseMiddleware, resty.ErrorHook) {
	afterFn := func(cli *resty.Client, response *resty.Response) error {
		logAccess(name, config, response.Request, response, nil)
		return nil
	}
	errorFn := func(req *resty.Request, err error) {
		if v, ok := err.(*resty.ResponseError); ok {
			logAccess(name, config, req, v.Response, v.Err)
		} else {
			logAccess(name, config, req, nil, err)
		}
	}
	return nil, afterFn, errorFn
}

func MetricInterceptor(name, addr string, r MetricPathRewriter) (resty.RequestMiddleware, resty.ResponseMiddleware, resty.ErrorHook) {
	rewriter := NoopMetricPathRewriter
	if r != nil {
		rewriter = r
	}

	afterFn := func(cli *resty.Client, res *resty.Response) error {
		method := res.Request.Method
		path := rewriter(res.Request.RawRequest.URL.Path)
		ClientHandleCounter.WithLabelValues(name, method, path, addr, strconv.Itoa(res.StatusCode())).Inc()
		ClientHandleHistogram.WithLabelValues(name, method, path, addr).Observe(res.Time().Seconds())
		return nil
	}

	errorFn := func(req *resty.Request, err error) {
		method := req.Method
		var path string

		// OnBeforeRequest 有错误时, 拿不到 req.RawRequest
		u, err2 := url.Parse(req.URL)
		if err2 != nil {
			path = "invalidUrl"
		} else {
			path = rewriter(u.Path)
		}

		if v, ok := err.(*resty.ResponseError); ok {
			ClientHandleCounter.WithLabelValues(name, method, path, addr, strconv.Itoa(v.Response.StatusCode())).Inc()
		} else {
			ClientHandleCounter.WithLabelValues(name, method, path, addr, "unknown").Inc()
		}

		ClientHandleHistogram.WithLabelValues(name, method, path, addr).Observe(time.Since(beg(req.Context())).Seconds())
	}

	return nil, afterFn, errorFn
}

func metricInterceptor(name string, config *Config) (resty.RequestMiddleware, resty.ResponseMiddleware, resty.ErrorHook) {
	addr := strings.TrimRight(config.Addr, "/")
	return MetricInterceptor(name, addr, config.MetricPathRewriter)
}

// Deprecated: use otel http transport
func traceInterceptor(name string, config *Config) (resty.RequestMiddleware, resty.ResponseMiddleware, resty.ErrorHook) { //nolint
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

// func customKeysInterceptor(name string, config *Config) (resty.RequestMiddleware, resty.ResponseMiddleware, resty.ErrorHook) {
// 	beforeFn := func(cli *resty.Client, req *resty.Request) error {
// 		transport.CustomKeysMapPropagator.Inject(req.Context(), propagation.HeaderCarrier(req.Header))
// 		return nil
// 	}

// 	return beforeFn, nil, nil
// }
