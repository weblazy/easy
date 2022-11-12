package grpc_client_config

import (
	"time"

	"github.com/weblazy/easy/utils/etrace"
	"github.com/weblazy/easy/utils/grpc/grpc_client/interceptor"
	"google.golang.org/grpc"
	"google.golang.org/grpc/balancer/roundrobin"
	"google.golang.org/grpc/keepalive"
)

// Config fgrpc client config.
type Config struct {
	Name             string
	Debug            bool          // 是否开启调试，默认不开启, 开启可以打印请求日志
	Addr             string        // 连接地址，直连为 127.0.0.1:9090，服务发现为 nacos:///appname
	BalancerName     string        // 负载均衡方式，默认 round_robin
	DialTimeout      time.Duration // 连接超时，默认3s
	ReadTimeout      time.Duration // 读超时，默认1s
	SlowLogThreshold time.Duration // 慢日志记录的阈值，默认600ms
	EnableBlock      bool          // 是否开启阻塞，默认开启
	// EnableOfficialGrpcLog        bool          // 是否开启官方grpc日志，默认关闭 // blog 和 zap 类型不兼容, 没法做
	EnableWithInsecure           bool // 是否开启非安全传输，默认开启
	EnableMetricInterceptor      bool // 是否开启监控，默认开启
	EnableTraceInterceptor       bool // 是否开启链路追踪，默认开启
	EnableAppNameInterceptor     bool // 是否开启传递应用名，默认开启
	EnableTimeoutInterceptor     bool // 是否开启超时传递，默认开启
	EnableServiceConfig          bool // 是否开启服务配置，默认关闭
	EnableFailOnNonTempDialError bool
	MetricSuccessCodes           []string // metric 监控, 统一将此列表中的 biz code rewrite 成统一成功 code 20000, 默认为空不做操作
	LogConf                      *interceptor.LogConf
	// Deprecated: not affect anything
	EnableSkyWalking bool // 是否额外开启 skywalking, 默认不开启

	KeepAlive   *keepalive.ClientParameters
	DialOptions []grpc.DialOption
}

// DefaultConfig defines grpc client default configuration
// User should construct config base on DefaultConfig
func DefaultConfig() *Config {
	return &Config{
		BalancerName:                 roundrobin.Name,
		DialTimeout:                  time.Second * 3,
		ReadTimeout:                  time.Second * 3,        // 3s
		SlowLogThreshold:             time.Millisecond * 600, // 600ms
		EnableBlock:                  true,
		EnableTraceInterceptor:       true,
		EnableWithInsecure:           true,
		EnableAppNameInterceptor:     true,
		EnableTimeoutInterceptor:     true,
		EnableMetricInterceptor:      true,
		EnableFailOnNonTempDialError: true,
		LogConf: &interceptor.LogConf{
			EnableAccessInterceptor:    true,
			EnableAccessInterceptorReq: true,
			EnableAccessInterceptorRes: true,
		},
		EnableServiceConfig: false,
		Addr:                "127.0.0.1:9090",
	}
}

func (config *Config) BuildDialOptions() {
	config.DialOptions = append(config.DialOptions, interceptor.GrpcHeaderCarrierInterceptor())

	// 最先执行trace
	if config.EnableTraceInterceptor {
		// 默认会启用 jaeger
		config.DialOptions = append(config.DialOptions, grpc.WithChainUnaryInterceptor(etrace.UnaryClientInterceptor()))
	}

	// 其次执行，自定义header头，这样才能赋值到ctx里
	// options = append(options, WithDialOption(grpc.WithChainUnaryInterceptor(customHeader(transport.CustomContextKeys()))))

	// 默认日志
	config.DialOptions = append(config.DialOptions, grpc.WithChainUnaryInterceptor(interceptor.LoggerUnaryClientInterceptor(config.LogConf)))

	if config.EnableTimeoutInterceptor {
		config.DialOptions = append(config.DialOptions, grpc.WithChainUnaryInterceptor(interceptor.TimeoutInterceptor(config.ReadTimeout)))
	}

	if config.EnableMetricInterceptor {
		config.DialOptions = append(config.DialOptions,
			grpc.WithChainUnaryInterceptor(interceptor.MetricUnaryClientInterceptor(config.MetricSuccessCodes)),
		)
	}

}
