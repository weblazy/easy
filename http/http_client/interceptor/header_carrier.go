package interceptor

import (
	"github.com/go-resty/resty/v2"
	"github.com/weblazy/easy/transport"
	"go.opentelemetry.io/otel/propagation"
)

// 多个服务间透传参数
func HeaderCarrierInterceptor() (resty.RequestMiddleware, resty.ResponseMiddleware, resty.ErrorHook) {
	beforeFn := func(cli *resty.Client, req *resty.Request) error {
		transport.CustomKeysMapPropagator.Inject(req.Context(), propagation.HeaderCarrier(req.Header))
		return nil
	}

	return beforeFn, nil, nil
}
