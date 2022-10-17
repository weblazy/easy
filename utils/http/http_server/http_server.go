package http_server

import (
	"context"
	"fmt"

	"github.com/fvbock/endless"
	"github.com/gin-gonic/gin"
	"github.com/weblazy/easy/utils/http/http_server/config"
	"github.com/weblazy/easy/utils/http/http_server/interceptor"
)

type HttpServer struct {
	Config *config.Config
	Engine *gin.Engine
}

func NewHttpServer(c *config.Config) (*HttpServer, error) {
	if c == nil {
		c = config.DefaultConfig()
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
	r.Use(interceptor.HeaderCarrierInterceptor())
	if server.Config.EnableTraceInterceptor {
		r.Use(interceptor.Trace(ctx))
	}
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
	return endless.ListenAndServe(fmt.Sprintf("%s:%d", s.Config.Host, s.Config.Port), s.Engine)
}
