package http_server_config

import (
	"time"

	"github.com/spf13/viper"
)

const PkgName = "http_server"

type Config struct {
	Name string
	Host string // IP地址，默认0.0.0.0
	Port int    // Port端口，默认80

	Timeout          time.Duration
	SlowLogThreshold time.Duration // 慢日志记录的阈值，默认 1s

	EnableTraceInterceptor  bool
	EnableMetricInterceptor bool
	EnableLogInterceptor    bool
	EnableAccessInterceptor bool // 是否开启记录请求数据，默认开启

}

// DefaultConfig default config ...
func DefaultConfig() *Config {
	return &Config{
		Host:                    "0.0.0.0",
		Port:                    80,
		Timeout:                 3 * time.Second,
		SlowLogThreshold:        time.Second,
		EnableTraceInterceptor:  true,
		EnableMetricInterceptor: true,
		EnableLogInterceptor:    true,
		EnableAccessInterceptor: true,
	}
}

func GetViperConfig(key string, cfg *viper.Viper) (*Config, error) {
	c := DefaultConfig()
	err := cfg.UnmarshalKey(key, c)
	if err != nil {
		return nil, err
	}
	return c, nil
}
