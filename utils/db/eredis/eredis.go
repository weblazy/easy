package eredis

import (
	"context"

	"github.com/go-redis/redis/v8"
	"github.com/weblazy/easy/utils/db/eredis/eredis_config"
	"github.com/weblazy/easy/utils/db/eredis/interceptor"

	"github.com/weblazy/easy/utils/elog"
)

var emptyCtx = context.Background()

type RedisClient struct {
	Config *eredis_config.Config
	redis.UniversalClient
}
type Option func(c *eredis_config.Config)

func NewRedisClient(c *eredis_config.Config, options ...Option) *RedisClient {
	if c == nil {
		c = eredis_config.DefaultConfig()
	}
	var client redis.UniversalClient
	switch c.Mode {
	case eredis_config.StubMode:
		client = NewClient(c)
	case eredis_config.SentinelMode:
		client = NewFailoverClient(c)
	case eredis_config.ClusterMode:
		client = NewClusterClient(c)
	default:
		panic(`redis mode must be one of ("stub", "cluster", "sentinel")`)
	}

	client.AddHook(interceptor.StartTimeHook())

	if c.EnableMetricInterceptor {
		client.AddHook(interceptor.MetricHook(c))
	}
	if c.EnableAccessInterceptor {
		client.AddHook(interceptor.LogHook(c))
	}

	if c.EnableTraceInterceptor {
		client.AddHook(interceptor.NewTracingHook())
	}
	for _, hook := range c.Hooks {
		client.AddHook(hook)
	}
	if err := client.Ping(emptyCtx).Err(); err != nil {
		switch c.OnFail {
		case eredis_config.OnFailPanic:
			elog.ErrorCtx(emptyCtx, "start redis panic", elog.FieldError(err))
		default:
			elog.ErrorCtx(emptyCtx, "start redis error", elog.FieldError(err))
		}
	}
	return &RedisClient{
		UniversalClient: client,
		Config:          c,
	}
}

// Client returns a universal redis client(ClusterClient, StubClient or SentinelClient), it depends on you config.
func (r *RedisClient) GetClient() redis.UniversalClient {
	return r
}

// Cluster try to get a redis.ClusterClient
func (r *RedisClient) GetCluster() *redis.ClusterClient {
	if c, ok := r.UniversalClient.(*redis.ClusterClient); ok {
		return c
	}
	return nil
}

// Stub try to get a redis.client
func (r *RedisClient) GetStub() *redis.Client {
	if c, ok := r.UniversalClient.(*redis.Client); ok {
		return c
	}
	return nil
}

// Sentinel try to get a redis Failover Sentinel client
func (r *RedisClient) GetSentinel() *redis.Client {
	if c, ok := r.UniversalClient.(*redis.Client); ok {
		return c
	}
	return nil
}

func NewClusterClient(c *eredis_config.Config) redis.UniversalClient {
	if len(c.Addrs) == 0 {
		panic(`invalid "addrs" config, "addrs" has none addresses but with cluster mode"`)
	}
	clusterClient := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:        c.Addrs,
		MaxRedirects: c.MaxRetries,
		ReadOnly:     c.ReadOnly,
		Password:     c.Password,
		MaxRetries:   c.MaxRetries,
		DialTimeout:  c.DialTimeout,
		ReadTimeout:  c.ReadTimeout,
		WriteTimeout: c.WriteTimeout,
		PoolSize:     c.PoolSize,
		MinIdleConns: c.MinIdleConns,
		IdleTimeout:  c.IdleTimeout,
	})
	return clusterClient
}

func NewClient(c *eredis_config.Config) redis.UniversalClient {
	if c.Addr == "" {
		panic(`invalid "addr" config, "addr" is empty but with stub mode"`)
	}
	client := redis.NewClient(&redis.Options{
		Addr:         c.Addr,
		Password:     c.Password,
		DB:           c.DB,
		MaxRetries:   c.MaxRetries,
		DialTimeout:  c.DialTimeout,
		ReadTimeout:  c.ReadTimeout,
		WriteTimeout: c.WriteTimeout,
		PoolSize:     c.PoolSize,
		MinIdleConns: c.MinIdleConns,
		IdleTimeout:  c.IdleTimeout,
	})
	return client
}

func NewFailoverClient(c *eredis_config.Config) redis.UniversalClient {
	if len(c.Addrs) == 0 {
		panic(`invalid "addrs" config, "addrs" has none addresses but with sentinel mode"`)
	}
	if c.MasterName == "" {
		panic(`invalid "masterName" config, "masterName" is empty but with sentinel mode"`)
	}
	failoverClient := redis.NewFailoverClient(&redis.FailoverOptions{
		MasterName:    c.MasterName,
		SentinelAddrs: c.Addrs,
		Password:      c.Password,
		DB:            c.DB,
		MaxRetries:    c.MaxRetries,
		DialTimeout:   c.DialTimeout,
		ReadTimeout:   c.ReadTimeout,
		WriteTimeout:  c.WriteTimeout,
		PoolSize:      c.PoolSize,
		MinIdleConns:  c.MinIdleConns,
		IdleTimeout:   c.IdleTimeout,
	})
	return failoverClient
}
