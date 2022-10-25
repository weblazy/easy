package eredis_config

import (
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
)

const (
	// ClusterMode using clusterClient
	ClusterMode string = "cluster"
	// SimpleMode using Client
	SimpleMode string = "simple"
	// FailoverMode using Failover sentinel client
	FailoverMode string = "failover"

	PkgName = "eredis"
)

// Config for redis, contains RedisStubConfig, RedisClusterConfig and RedisSentinelConfig
type Config struct {
	Name       string   // Name redis名称
	Addrs      []string // Addrs Cluster,Failover实例配置地址
	Addr       string   // Addr Simple 实例配置地址
	Mode       string   // Mode Redis模式 cluster|simple|failover
	MasterName string   // MasterName 哨兵主节点名称，sentinel模式下需要配置此项
	Password   string   // Password 密码
	DB         int      // DB，默认为0, 一般应用不推荐使用DB分片
	PoolSize   int      // PoolSize 集群内每个节点的最大连接池限制 默认每个CPU10个连接

	MaxRetries   int           // MaxRetries 网络相关的错误最大重试次数 默认8次
	MinIdleConns int           // MinIdleConns 最小空闲连接数
	DialTimeout  time.Duration // DialTimeout 拨超时时间
	ReadTimeout  time.Duration // ReadTimeout 读超时 默认3s
	WriteTimeout time.Duration // WriteTimeout 读超时 默认3s
	IdleTimeout  time.Duration // IdleTimeout 连接最大空闲时间，默认60s, 超过该时间，连接会被主动关闭
	ReadOnly     bool          // ReadOnly 集群模式 在从属节点上启用读模式

	EnableMetricInterceptor bool // 是否开启监控，默认开启
	EnableTraceInterceptor  bool // 是否开启链路，默认开启

	SlowLogThreshold time.Duration // 慢日志门限值，超过该门限值的请求，将被记录到慢日志中
	EnableLogAccess  bool          // 是否开启，成功时也记录请求日志
	EnableLogReq     bool          // 是否开启记录请求参数
	EnableLogRes     bool          // 是否开启记录响应参数
	Hooks            []redis.Hook
}

// DefaultConfig default config ...
func DefaultConfig() *Config {
	return &Config{
		Mode:                    SimpleMode,
		DB:                      0,
		PoolSize:                0, // will be handled by redis v8
		MaxRetries:              0,
		MinIdleConns:            20,
		DialTimeout:             time.Second,
		ReadTimeout:             time.Second,
		WriteTimeout:            time.Second,
		IdleTimeout:             time.Second * 60,
		ReadOnly:                false,
		SlowLogThreshold:        time.Millisecond * 250,
		EnableMetricInterceptor: true,
		EnableTraceInterceptor:  true,
		EnableLogAccess:         false,
		EnableLogReq:            true,
		EnableLogRes:            true,
	}
}

// AddrString 获取地址, 用于监控
// 多个地址会用 , 连接
func (c Config) AddrString() string {
	addr := c.Addr
	if len(c.Addrs) > 0 {
		addr = strings.Join(c.Addrs, ",")
	}
	return addr
}
