package interceptor

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"runtime"
	"time"

	"github.com/weblazy/easy/utils/code_err"
	"github.com/weblazy/easy/utils/etrace"
	"github.com/weblazy/easy/utils/timex"

	"github.com/weblazy/easy/utils/http/http_server/http_server_config"
	"github.com/weblazy/easy/utils/http/http_server/service"

	"github.com/weblazy/easy/utils/elog"

	gocorezap "github.com/weblazy/easy/utils/elog/zap"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func init() {
	gocorezap.InitFileLog("log/access")
}

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
	return func(c *gin.Context) {
		if c.Request.Method == http.MethodGet {

		} else if c.ContentType() == gin.MIMEJSON {
			LogJson(c, cfg)
		} else if c.ContentType() == gin.MIMEMultipartPOSTForm {

		}

	}
}

func LogJson(c *gin.Context, cfg *http_server_config.Config) {
	req := c.Request
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
		duration := time.Since(GetStartTime(req.Context()))
		fields := []zap.Field{
			zap.String("url", req.URL.String()),
			zap.String("host", c.GetHeader("Host")),
			zap.String("path", req.URL.Path),
			zap.String("method", c.GetHeader("Method")),
			zap.Any("req_header", req.Header),
			zap.String("req_body", logData.RequestBody),
			zap.Any("res_header", c.Writer.Header()),
			zap.String("res_body", blw.body.String()),
			zap.String("client_ip", c.ClientIP()),
			zap.String("start_time", GetStartTime(req.Context()).Format(timex.TimeLayout)),
			zap.Float64("duration", float64(duration.Microseconds())/1000),
		}
		fields = append(fields, zap.Int("status_code", c.Writer.Status()))
		if cfg.SlowLogThreshold > time.Duration(0) && duration > cfg.SlowLogThreshold {
			fields = append(fields, zap.Bool("slow", true))
		}
		// 开启了链路，那么就记录链路id
		if cfg.EnableTraceInterceptor && etrace.IsGlobalTracerRegistered() {
			fields = append(fields, zap.String("trace_id", etrace.ExtractTraceID(req.Context())))
		}
		if err != nil {
			fields = append(fields, zap.String("event", "error"), zap.Error(err))
			elog.WarnCtx(req.Context(), "http_server", fields...)
			return
		} else if cfg.EnableAccessInterceptor {
			fields = append(fields, zap.String("event", "normal"))
			elog.InfoCtx(req.Context(), "http_server", fields...)
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
