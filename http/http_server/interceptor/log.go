package interceptor

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"runtime"
	"sync"
	"time"

	"github.com/weblazy/easy/code_err"
	"github.com/weblazy/easy/etrace"
	"github.com/weblazy/easy/timex"

	"github.com/weblazy/easy/http/http_server/http_server_config"
	"github.com/weblazy/easy/http/http_server/service"

	"github.com/weblazy/easy/elog"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

var once sync.Once
var emptyData = struct{}{}

type LogData struct {
	RequestBody string `json:"request_body"`
}

type BodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w BodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func (w BodyLogWriter) WriteString(s string) (int, error) {
	w.body.WriteString(s)
	return w.ResponseWriter.WriteString(s)
}

// Log returns a middleware
// and handles the control to the centralized HTTPErrorHandler.
func Log(ctx context.Context, cfg *http_server_config.Config) gin.HandlerFunc {
	once.Do(cfg.InitLogger)
	return func(c *gin.Context) {
		if c.ContentType() == gin.MIMEJSON {
			LogJson(c, cfg)
		} else if c.ContentType() == gin.MIMEMultipartPOSTForm {

		}

	}
}

func LogJson(c *gin.Context, cfg *http_server_config.Config) {
	req := c.Request
	ctx := elog.SetLogerName(req.Context(), http_server_config.PkgName)
	logData := &LogData{}
	blw := &BodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
	c.Writer = blw
	var bodyBytes []byte
	bodyBytes, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		Error(c, code_err.ParamsErr, fmt.Errorf("Invalid request body"))
		return
	}
	logData.RequestBody = string(bodyBytes)
	defer func() {
		duration := time.Since(GetStartTime(ctx))
		fields := []zap.Field{
			zap.String("url", req.URL.String()),
			zap.String("host", req.Host),
			zap.String("path", req.URL.Path),
			elog.FieldMethod(req.Method),
			zap.Any("req_header", req.Header),
			zap.String("req_body", logData.RequestBody),
			zap.Any("res_header", c.Writer.Header()),
			zap.String("res_body", blw.body.String()),
			zap.String("client_ip", c.ClientIP()),
			zap.String("start_time", GetStartTime(ctx).Format(timex.TimeLayout)),
			elog.FieldDuration(duration),
		}
		fields = append(fields, zap.Int("status_code", c.Writer.Status()))
		var isSlow bool
		if cfg.SlowLogThreshold > time.Duration(0) && duration > cfg.SlowLogThreshold {
			isSlow = true
		}
		fields = append(fields, elog.FieldSlow(isSlow))
		// 开启了链路，那么就记录链路id
		if cfg.EnableTraceInterceptor && etrace.IsGlobalTracerRegistered() {
			fields = append(fields, elog.FieldTrace(etrace.ExtractTraceID(ctx)))
		}
		if err != nil {
			fields = append(fields, zap.String("event", "error"), zap.Error(err))
			elog.ErrorCtx(ctx, http_server_config.PkgName, fields...)
			return
		} else if isSlow {
			elog.WarnCtx(ctx, http_server_config.PkgName, fields...)
		} else if cfg.EnableAccessInterceptor {
			fields = append(fields, zap.String("event", "normal"))
			elog.InfoCtx(ctx, http_server_config.PkgName, fields...)
		}
	}()

	// 新建缓冲区并替换原有Request.body
	c.Request.Body = ioutil.NopCloser(bytes.NewBuffer([]byte(logData.RequestBody)))
	c.Next()
}

func Error(c *gin.Context, codeErr *code_err.CodeErr, err error) {
	_, file, line, ok := runtime.Caller(1)
	if ok {
		elog.InfoCtx(c.Request.Context(), fmt.Sprintf("%s:%d %s", file, line, err.Error()))
	}
	resp := &service.Response{
		Code: codeErr.Code,
		Msg:  codeErr.Msg,
	}
	c.JSON(200, resp)
	c.Abort()
}
