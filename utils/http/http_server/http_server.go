package http_server

import (
	"context"
	"fmt"

	"github.com/fvbock/endless"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"

	"github.com/weblazy/easy/utils/http/http_server/http_server_config"
	"github.com/weblazy/easy/utils/http/http_server/interceptor"
)

type HttpServer struct {
	Config *http_server_config.Config
	*gin.Engine
}

func NewHttpServerViper(key string, cfg *viper.Viper) (*HttpServer, error) {
	c := http_server_config.DefaultConfig()
	cfg.UnmarshalKey(key, c)
	server := &HttpServer{
		Config: c,
	}
	return server, nil
}

func NewHttpServer(c *http_server_config.Config) (*HttpServer, error) {
	if c == nil {
		c = http_server_config.DefaultConfig()
	}

	server := &HttpServer{
		Config: c,
	}
	ctx := context.Background()
	// opts = append([]RunOption{WithNotFoundHandler(nil)}, opts...)
	// for _, opt := range opts {
	// 	opt(server)
	// }
	r := gin.New()
	r.Use(interceptor.SetStartTimeInterceptor())
	if server.Config.EnableTraceInterceptor {
		r.Use(otelgin.Middleware(c.Name))
		r.Use(interceptor.Trace(ctx))
	}
	r.Use(interceptor.HeaderCarrierInterceptor())

	if server.Config.EnableLogInterceptor {
		r.Use(interceptor.Log(ctx, c))
	}
	if server.Config.EnableMetricInterceptor {
		r.Use(interceptor.MetricInterceptor(c))
	}
	if server.Config.Timeout > 0 {
		r.Use(interceptor.Timeout(server.Config.Timeout))
	}
	r.Use(gin.Recovery())
	server.Engine = r
	return server, nil
}

func (s *HttpServer) Start() error {
	return endless.ListenAndServe(fmt.Sprintf("%s:%d", s.Config.Host, s.Config.Port), s)
}
