package ejager

import (
	"context"
	"os"

	"github.com/weblazy/easy/utils/glog"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

const (
	ProjectEnvKey = "projectEnv"
	EnvKey        = "env"
)

const (
	MY_ENV_NAME_KEY         = "BSM_SERVICE_STAGE"
	MY_PROJECT_ENV_NAME_KEY = "MY_PROJECT_ENV_NAME"
	MY_PROJECT_NAME_KEY     = "MY_PROJECT_NAME"
)

var emptyCtx = context.Background()

type Config struct {
	Enable bool

	ServiceName    string
	EnvName        string
	ProjectEnvName string

	AgentHost string // agent host
	AgentPort string // agent port
	Fraction  float64

	options []tracesdk.TracerProviderOption
}

func DefaultConfig() *Config {
	return &Config{
		Enable:         true,
		ServiceName:    os.Getenv(MY_PROJECT_NAME_KEY),
		EnvName:        os.Getenv(MY_ENV_NAME_KEY),
		ProjectEnvName: os.Getenv(MY_PROJECT_ENV_NAME_KEY),
		AgentHost:      "jaeger-agent-cluster.inner-udp.efficiency.ww5sawfyut0k.bitsvc.io",
		AgentPort:      "6831",
		Fraction:       1.0,
	}
}

// Option 可选项
type Option func(c *Config)

func (config *Config) WithTracerProviderOption(options ...tracesdk.TracerProviderOption) *Config {
	config.options = append(config.options, options...)
	return config
}

func (config *Config) Build(opts ...Option) trace.TracerProvider {
	lc := zap.Any("config", config)

	if !config.Enable {
		glog.InfoCtx(emptyCtx, "jaeger not enable", lc)
		return trace.NewNoopTracerProvider()
	}

	if config.ServiceName == "" {
		glog.InfoCtx(emptyCtx, "jaeger not enable, empty ServiceName", lc)
		return trace.NewNoopTracerProvider()
	}

	endpoint := jaeger.WithAgentEndpoint(jaeger.WithAgentHost(config.AgentHost), jaeger.WithAgentPort(config.AgentPort))
	exp, err := jaeger.New(endpoint)

	if err != nil {
		glog.InfoCtx(emptyCtx, "init jaeger client error", lc, zap.Error(err))
		return trace.NewNoopTracerProvider()
	}

	options := []tracesdk.TracerProviderOption{
		// Set the sampling rate based on the parent span to 100%
		// otel 暂未实现
		// 但是我们基本不是 span 采样的源头, 所以基本只需要集成 parent 就行
		// https://github.com/open-telemetry/opentelemetry-go-contrib/pull/936
		tracesdk.WithSampler(tracesdk.ParentBased(tracesdk.TraceIDRatioBased(config.Fraction))),
		// Always be sure to batch in production.
		tracesdk.WithBatcher(exp),
		// Record information about this application in a Resource.
		tracesdk.WithResource(resource.NewSchemaless(
			semconv.ServiceNameKey.String(config.ServiceName),
			semconv.DeploymentEnvironmentKey.String(config.EnvName),
			attribute.String(EnvKey, config.EnvName),
			attribute.Key(ProjectEnvKey).String(config.ProjectEnvName),
		)),
	}

	for _, opt := range opts {
		opt(config)
	}

	options = append(options, config.options...)
	tp := tracesdk.NewTracerProvider(options...)

	glog.InfoCtx(emptyCtx, "jaeger init success", lc)

	return tp
}
