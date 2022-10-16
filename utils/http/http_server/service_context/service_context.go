package service_context

import (
	"context"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sunmi-OS/gocore/v2/utils"
)

type ServiceContext struct {
	*gin.Context
	Ctx context.Context
	R   Response
}

var (
	ErrorBind = errors.New("missing required parameters")
)

//const TraceHeaderKey = "TraceHeaderKey"

// NewContext 初始化上下文包含context.Context
// 对链路信息进行判断并且在Response时返回TraceId信息
func NewServiceContext(g *gin.Context) ServiceContext {
	c := ServiceContext{
		Context: g,
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
	// c.R.Code = ecode.Transform(err)
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
