package interceptor

import (
	"github.com/gin-gonic/gin"
	"github.com/weblazy/easy/transport"
	"go.opentelemetry.io/otel/propagation"
)

// 多个服务间透传参数
func HeaderCarrierInterceptor() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Request = c.Request.WithContext(transport.CustomKeysMapPropagator.Extract(c.Request.Context(), propagation.HeaderCarrier(c.Request.Header)))
	}
}
