package grpc_server_config

import (
	"fmt"
	"time"

	"github.com/weblazy/easy/utils/elog"
	"github.com/weblazy/easy/utils/elog/ezap"
	"google.golang.org/grpc"
)

const (
	DefaultPort = 9090
	PkgName     = "grpc_server"
)

// Config ...
type Config struct {
	Name                       string
	Host                       string        // IP地址，默认0.0.0.0
	Port                       int           // Port端口，默认9090
	Network                    string        // 网络类型，默认tcp4
	EnableMetricInterceptor    bool          // 是否开启监控，默认开启
	EnableTraceInterceptor     bool          // 是否开启链路追踪，默认开启
	EnableSkipHealthLog        bool          // 是否屏蔽探活日志，默认关闭
	SlowLogThreshold           time.Duration // 服务慢日志，默认500ms
	EnableAccessInterceptor    bool          // 是否开启，记录请求数据
	EnableAccessInterceptorReq bool          // 是否开启记录请求参数，默认不开启
	EnableAccessInterceptorRes bool          // 是否开启记录响应参数，默认不开启
	EnableServerReflection     bool          // 是否开启 reflection, 默认开启
	EnableHealth               bool          // 是否开启 grpc health, 默认开启
	MinDeadlineDuration        time.Duration // server handler ctx 最短超时时间, 默认 10s
	MetricSuccessCodes         []string      // metric 监控, 统一将此列表中的 biz code rewrite 成统一成功 code 20000, 默认为空不做操作
	// Deprecated: not affect anything
	EnableSkyWalking bool // 是否额外开启 skywalking, 默认开启

	ServerOptions            []grpc.ServerOption
	StreamInterceptors       []grpc.StreamServerInterceptor
	UnaryInterceptors        []grpc.UnaryServerInterceptor
	PrependUnaryInterceptors []grpc.UnaryServerInterceptor

	EnableFielLogger bool // 将日志输出到文件
	FielLoggerPath   string
}

// DefaultConfig represents default config
// User should construct config base on DefaultConfig
func DefaultConfig() *Config {
	return &Config{
		Network:                    "tcp4",
		Host:                       "0.0.0.0",
		Port:                       DefaultPort,
		EnableMetricInterceptor:    true,
		EnableSkipHealthLog:        true,
		EnableTraceInterceptor:     true,
		SlowLogThreshold:           time.Millisecond * 800, // 800ms
		EnableAccessInterceptor:    true,
		EnableAccessInterceptorReq: true,
		EnableAccessInterceptorRes: true,
		EnableServerReflection:     true,
		EnableHealth:               true,
		MinDeadlineDuration:        time.Second * 10,
		ServerOptions:              []grpc.ServerOption{},
		StreamInterceptors:         []grpc.StreamServerInterceptor{},
		UnaryInterceptors:          []grpc.UnaryServerInterceptor{},
		FielLoggerPath:             PkgName,
	}
}

// Address ...
func (config Config) Address() string {
	return fmt.Sprintf("%s:%d", config.Host, config.Port)
}

func (config Config) InitLogger() {
	if config.EnableFielLogger {
		logger := ezap.NewFileEzap(config.FielLoggerPath)
		elog.SetLogger(PkgName, logger)
		return
	}
	elog.SetLogger(PkgName, elog.DefaultLogger)
}
