package service

import (
	"context"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sunmi-OS/gocore/v2/utils"
	"github.com/weblazy/easy/utils/code_err"
	"go.uber.org/zap"
)

type ServiceContext struct {
	*gin.Context
	Ctx context.Context
	R   Response
}

var (
	ErrorBind            = errors.New("missing required parameters")
	defaultErrCode int64 = -1
)

// NewContext 初始化上下文包含context.Context
func NewServiceContext(g *gin.Context) ServiceContext {
	c := ServiceContext{
		Context: g,
		Ctx:     g.Request.Context(),
		R:       NewResponse(),
	}

	return c
}

// Success 返回正常数据
func (c *ServiceContext) Success(data interface{}) {
	c.R.Data = data
	c.JSON(http.StatusOK, c.R)
}

// Error 返回异常信息，自动识别Code码
func (c *ServiceContext) Error(err error) {
	if e, ok := err.(*code_err.CodeErr); ok {
		c.R.Code = e.Code
	}
	if c.R.Code == 0 {
		c.R.Code = defaultErrCode
	}
	c.R.Msg = err.Error()
	c.JSON(http.StatusOK, c.R)
}

// ErrorCodeMsg 直接指定code和msg
func (c *ServiceContext) ErrorCodeMsg(code int64, msg string) {
	c.R.Code = code
	c.R.Msg = msg
	c.JSON(http.StatusOK, c.R)
}

// Response 直接指定code和msg和data
func (c *ServiceContext) Response(code int64, msg string, data interface{}) {
	c.R.Code = code
	c.R.Msg = msg
	c.R.Data = data
	c.JSON(http.StatusOK, c.R)
}

// 打印log
func (c *ServiceContext) Log(msg string, codeErr *code_err.CodeErr, err error) error {
	e := code_err.Log(c.Ctx, msg, codeErr, err)
	c.Error(e)
	return e
}

// 打印log
func (c *ServiceContext) ErrLog(codeErr *code_err.CodeErr, err error) error {
	e := code_err.ErrLog(c.Ctx, codeErr, err)
	c.Error(e)
	return e
}

// 打印log
func (c *ServiceContext) ErrLogf(codeErr *code_err.CodeErr, format string, a ...interface{}) error {
	e := code_err.ErrLogf(c.Ctx, codeErr, format, a...)
	c.Error(e)
	return e
}

// 打印log
func (c *ServiceContext) LogField(codeErr *code_err.CodeErr, msg string, fields ...zap.Field) error {
	e := code_err.LogField(c.Ctx, codeErr, msg, fields...)
	c.Error(e)
	return e
}

// BindValidator 参数绑定结构体，并且按照tag进行校验返回校验结果
func (c *ServiceContext) BindValidator(obj interface{}) error {
	err := c.ShouldBind(obj)
	if err != nil {
		if utils.IsRelease() {
			return ErrorBind
		}
		return err
	}
	return nil
}
