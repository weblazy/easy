package eredis

import (
	"context"

	"github.com/go-redis/redis/v8"
	"github.com/weblazy/easy/utils/elog"
)

const (
	OnFailPanic = "panic"
)

var emptyCtx = context.Background()

type RedisClient struct {
	Config *Config
	Client redis.UniversalClient
}

func NewRedisClient(c *Config) *RedisClient {
	if c == nil {
		c = DefaultConfig()
	}
	var client redis.UniversalClient
	switch c.Mode {
	case StubMode:
		client = NewClient(c)
	case SentinelMode:
		client = NewFailoverClient(c)
	case ClusterMode:
		client = NewClusterClient(c)
	default:
		panic(`redis mode must be one of ("stub", "cluster", "sentinel")`)
	}
	if c.EnableTraceInterceptor {
		client.AddHook(NewTracingHook())
	}
	for _, incpt := range c.interceptors {
		client.AddHook(incpt)
	}
	if err := client.Ping(emptyCtx).Err(); err != nil {
		switch c.OnFail {
		case OnFailPanic:
			elog.ErrorCtx(emptyCtx, "start redis panic", elog.FieldError(err))
		default:
			elog.ErrorCtx(emptyCtx, "start redis error", elog.FieldError(err))
		}
	}
	return &RedisClient{
		Config: c,
		Client: client,
	}
}

// Client returns a universal redis client(ClusterClient, StubClient or SentinelClient), it depends on you config.
func (r *RedisClient) GetClient() redis.UniversalClient {
	return r.Client
}

// Cluster try to get a redis.ClusterClient
func (r *RedisClient) GetCluster() *redis.ClusterClient {
	if c, ok := r.Client.(*redis.ClusterClient); ok {
		return c
	}
	return nil
}

// Stub try to get a redis.client
func (r *RedisClient) GetStub() *redis.Client {
	if c, ok := r.Client.(*redis.Client); ok {
		return c
	}
	return nil
}

// Sentinel try to get a redis Failover Sentinel client
func (r *RedisClient) GetSentinel() *redis.Client {
	if c, ok := r.Client.(*redis.Client); ok {
		return c
	}
	return nil
}

func NewClusterClient(c *Config) redis.UniversalClient {
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

func NewClient(c *Config) redis.UniversalClient {
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

func NewFailoverClient(c *Config) redis.UniversalClient {
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
