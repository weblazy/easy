package interceptor

import (
	"net/url"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/weblazy/easy/utils/etrace"
	"github.com/weblazy/easy/utils/glog"
	"github.com/weblazy/easy/utils/http/http_client/config"
	"go.uber.org/zap"
)

func LogInterceptor(cfg *config.Config) (resty.RequestMiddleware, resty.ResponseMiddleware, resty.ErrorHook) {
	afterFn := func(cli *resty.Client, response *resty.Response) error {
		logAccess(cfg, response.Request, response, nil)
		return nil
	}
	errorFn := func(req *resty.Request, err error) {
		if v, ok := err.(*resty.ResponseError); ok {
			logAccess(cfg, req, v.Response, v.Err)
		} else {
			logAccess(cfg, req, nil, err)
		}
	}
	return nil, afterFn, errorFn
}

func logAccess(cfg *config.Config, req *resty.Request, res *resty.Response, err error) {
	rawRequest := req.RawRequest
	var path, host string
	// 修复err 不是 *resty.ResponseError错误的时候，可能为nil
	if rawRequest != nil {
		path = rawRequest.URL.RequestURI()
		host = rawRequest.URL.Host
	} else { // RawRequest 不一定总是有
		u, err2 := url.Parse(req.URL)
		if err2 == nil {
			path = u.RequestURI()
			host = u.Host
		}
	}

	var duration = time.Since(GetStartTime(req.Context()))
	var respBody string
	if res != nil {
		respBody = string(res.Body())
	}

	var fields = make([]zap.Field, 0, 20)
	fields = append(fields,
		zap.String("name", cfg.Name),
		zap.String("host", host),
		zap.String("method", req.Method),
		zap.String("path", path),
		zap.Float64("duration", float64(duration.Microseconds())/1000),
	)

	// 开启了链路，那么就记录链路id
	if cfg.EnableTraceInterceptor && etrace.IsGlobalTracerRegistered() {
		fields = append(fields, zap.String("trace_id", etrace.ExtractTraceID(req.Context())))
	}

	if cfg.EnableAccessInterceptorReq {
		if cfg.EnableAccessInterceptorReqHeader {
			fields = append(fields, zap.Any("req_header", req.Header))
		}
		fields = append(fields, zap.Any("req_body", req.Body))
	}

	if cfg.EnableAccessInterceptorRes {
		// 处理 res 为空时空指针错误
		if res != nil {
			fields = append(fields, zap.Any("res_header", res.Header()))
			fields = append(fields, zap.Any("status_code", res.StatusCode()))
		}
		fields = append(fields, zap.Any("res_body", respBody))
	}

	if cfg.SlowLogThreshold > time.Duration(0) && duration > cfg.SlowLogThreshold {
		fields = append(fields, zap.Bool("slow", true))
	}

	if err != nil {
		fields = append(fields, zap.String("event", "error"), zap.Error(err))
		if res == nil {
			// 无 res 的是连接超时等系统级错误
			glog.ErrorCtx(req.Context(), "access", fields...)
			return
		}
		glog.WarnCtx(req.Context(), "access", fields...)
		return
	}

	if cfg.EnableAccessInterceptor {
		fields = append(fields, zap.String("event", "normal"))
		glog.InfoCtx(req.Context(), "access", fields...)
	}
}