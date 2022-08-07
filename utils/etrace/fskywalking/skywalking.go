package fskywalking

import (
	"context"
	"fmt"
	"os"

	"go.uber.org/zap"

	"github.com/SkyAPM/go2sky"
	"github.com/SkyAPM/go2sky/reporter"
	"github.com/weblazy/easy/utils/glog"
)

const (
	ENV_KEY          = "MY_ENV_NAME"
	PROJECT_NAME_KEY = "MY_PROJECT_NAME"
)

var emptyCtx = context.Background()

type Config struct {
	Enable      bool
	ServiceName string
	EnvName     string

	AgentEndPoint string
	Sampler       float64
}

func DefaultConfig() *Config {
	return &Config{
		Enable:        true,
		ServiceName:   os.Getenv(PROJECT_NAME_KEY),
		EnvName:       os.Getenv(ENV_KEY),
		AgentEndPoint: "monitor.infra.ww5sawfyut0k.bitsvc.io:31801",
		Sampler:       0.1,
	}
}

// Option 可选项
type Option func(c *Config)

func (config *Config) Build(ops ...Option) *go2sky.Tracer {
	lc := zap.Any("config", config)

	if !config.Enable {
		glog.InfoCtx(emptyCtx, "skywalking not enable", lc)
		return nil
	}

	r, err := reporter.NewGRPCReporter(config.AgentEndPoint)
	if err != nil {
		glog.InfoCtx(emptyCtx, "skywalking new reporter error", lc, zap.Error(err))
		return nil
	}

	tracer, err := go2sky.NewTracer(fmt.Sprintf("%s-%s", config.EnvName, config.ServiceName), go2sky.WithReporter(r), go2sky.WithSampler(config.Sampler))
	if err != nil {
		glog.InfoCtx(emptyCtx, "skywalking new tracer error", lc, zap.Error(err))
		return nil
	}

	// registers `tracer` as the global Tracer
	go2sky.SetGlobalTracer(tracer)

	glog.InfoCtx(emptyCtx, "skywalking init success", lc)

	return tracer
}