package ehttp

import (
	"runtime"
	"time"
)

// Config HTTP配置选项
type Config struct {
	Addr                             string             // 连接地址
	Debug                            bool               // 是否开启调试，默认不开启，开启后并加上export EGO_DEBUG=true，可以看到每次请求，配置名、地址、耗时、请求数据、响应数据
	RawDebug                         bool               // 是否开启原生调试，默认不开启
	ReadTimeout                      time.Duration      // 读超时，默认 3s
	SlowLogThreshold                 time.Duration      // 慢日志记录的阈值，默认 1s
	IdleConnTimeout                  time.Duration      // 设置空闲连接时间，默认90 * time.Second
	MaxIdleConns                     int                // 设置最大空闲连接数
	MaxIdleConnsPerHost              int                // 设置长连接个数
	EnableMetricInterceptor          bool               // 是否开启 metric, 默认关闭
	EnableTraceInterceptor           bool               // 是否开启链路追踪，默认开启
	EnableKeepAlives                 bool               // 是否开启长连接，默认打开
	EnableAccessInterceptor          bool               // 是否开启记录请求数据，默认开启
	EnableAccessInterceptorReq       bool               // 是否开启记录请求参数，默认开启
	EnableAccessInterceptorReqHeader bool               // 是否开启记录请求 header 参数，默认关闭
	EnableAccessInterceptorRes       bool               // 是否开启记录响应参数，默认开启
	Proxy                            string             // 支持配置显示传递代理，如：http://
	MetricPathRewriter               MetricPathRewriter // 指标监控 path 重写方法, 防止 metrics label 不可控
}

// DefaultConfig ...
func DefaultConfig() *Config {
	return &Config{
		Debug:                      false,
		RawDebug:                   false,
		ReadTimeout:                time.Second * 3,
		SlowLogThreshold:           time.Second,
		MaxIdleConns:               100,
		MaxIdleConnsPerHost:        runtime.GOMAXPROCS(0) + 1,
		IdleConnTimeout:            90 * time.Second,
		EnableKeepAlives:           true,
		EnableTraceInterceptor:     true,
		EnableAccessInterceptor:    true,
		EnableAccessInterceptorReq: true,
		EnableAccessInterceptorRes: true,
		MetricPathRewriter:         NoopMetricPathRewriter,
	}
}
