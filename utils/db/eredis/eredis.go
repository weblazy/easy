package eredis

import (
	"context"
	"fmt"

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

func NewRedisClient(c *eredis_config.Config) *RedisClient {
	if c == nil {
		c = eredis_config.DefaultConfig()
	}
	var client redis.UniversalClient
	switch c.Mode {
	case eredis_config.SimpleMode:
		client = NewClient(c)
	case eredis_config.FailoverMode:
		client = NewFailoverClient(c)
	case eredis_config.ClusterMode:
		client = NewClusterClient(c)
	default:
		elog.ErrorCtx(emptyCtx, "redis mode must be one of (simple, cluster, failover)", elog.FieldName(c.Name))
		return nil
	}

	client.AddHook(interceptor.StartTimeHook())

	if c.EnableTraceInterceptor {
		client.AddHook(interceptor.NewTracingHook())
	}
	if c.EnableMetricInterceptor {
		client.AddHook(interceptor.MetricHook(c))
	}
	client.AddHook(interceptor.LogHook(c))
	for _, hook := range c.Hooks {
		client.AddHook(hook)
	}
	if err := client.Ping(emptyCtx).Err(); err != nil {
		elog.ErrorCtx(emptyCtx, "start redis error", elog.FieldName(c.Name), elog.FieldError(err))
		return nil
	}
	return &RedisClient{
		UniversalClient: client,
		Config:          c,
	}
}

// GetUniversalClient returns a universal redis client(ClusterClient, SimpleClient or FailoverClient), it depends on you config.
func (r *RedisClient) GetUniversalClient() redis.UniversalClient {
	return r
}

// GetClient try to get a redis.client
func (r *RedisClient) GetClient() *redis.Client {
	if c, ok := r.UniversalClient.(*redis.Client); ok {
		return c
	}
	return nil
}

// GetClusterClient try to get a redis.ClusterClient
func (r *RedisClient) GetClusterClient() *redis.ClusterClient {
	if c, ok := r.UniversalClient.(*redis.ClusterClient); ok {
		return c
	}
	return nil
}

// GetFailoverClient try to get a redis Failover Sentinel client
func (r *RedisClient) GetFailoverClient() *redis.Client {
	if c, ok := r.UniversalClient.(*redis.Client); ok {
		return c
	}
	return nil
}

func NewClient(c *eredis_config.Config) redis.UniversalClient {
	if c.Addr == "" {
		elog.ErrorCtx(emptyCtx, eredis_config.PkgName, elog.FieldName(c.Name), elog.FieldError(fmt.Errorf(`invalid "addr" config, "addr" is empty but with stub mode"`)))
		return nil
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

func NewClusterClient(c *eredis_config.Config) redis.UniversalClient {
	if len(c.Addrs) == 0 {
		elog.ErrorCtx(emptyCtx, "invalid addrs config, addrs has none addresses but with cluster mode", elog.FieldName(c.Name))
		return nil
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

func NewFailoverClient(c *eredis_config.Config) redis.UniversalClient {
	if len(c.Addrs) == 0 {
		elog.ErrorCtx(emptyCtx, `invalid "addrs" config, "addrs" has none addresses but with failover mode`, elog.FieldName(c.Name))
		return nil
	}
	if c.MasterName == "" {
		elog.ErrorCtx(emptyCtx, `invalid "masterName" config, "masterName" is empty but with sentinel mode"`, elog.FieldName(c.Name))
		return nil
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
