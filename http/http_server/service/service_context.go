package service

import (
	"context"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/weblazy/easy/code_err"
	"github.com/weblazy/easy/env"
)

type ServiceContext struct {
	*gin.Context
	*code_err.Log
	R   Response
	Ctx context.Context //gin.Context.Request.Context
}

var (
	ErrorBind            = errors.New("missing required parameters")
	defaultErrCode int64 = -1
)

// NewContext 初始化上下文包含context.Context
func NewServiceContext(g *gin.Context) *ServiceContext {
	ctx := g.Request.Context()
	c := ServiceContext{
		Context: g,
		Log:     code_err.NewLog(ctx),
		R:       NewResponse(),
		Ctx:     ctx,
	}

	return &c
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
func (c *ServiceContext) Return(err *code_err.CodeErr) {
	if err != nil {
		c.R.Code = err.Code
		c.R.Msg = err.Msg
	}
	c.JSON(http.StatusOK, c.R)
}

// Success 返回正常数据
func (c *ServiceContext) SetData(data interface{}) *code_err.CodeErr {
	c.R.Data = data
	return nil
}

// BindValidator 参数绑定结构体，并且按照tag进行校验返回校验结果
func (c *ServiceContext) BindValidator(obj interface{}) error {
	err := c.ShouldBind(obj)
	if err != nil {
		if env.IsRelease() {
			return ErrorBind
		}
		return err
	}
	return nil
}
