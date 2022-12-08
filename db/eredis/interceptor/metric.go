package interceptor

import (
	"context"
	"errors"

	"github.com/go-redis/redis/v8"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/weblazy/easy/db/eredis/eredis_config"
)

var (

	// ClientHandleCounter ...
	RedisHandleCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "",
			Name:      "redis_handle_total",
		}, []string{"name", "method", "addr", "result"})

	// ClientHandleHistogram ...
	RedisHandleHistogram = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "",
		Name:      "redis_handle_seconds",
	}, []string{"name", "method", "addr"})
)

func init() {
	prometheus.MustRegister(RedisHandleCounter)
	prometheus.MustRegister(RedisHandleHistogram)
}

func MetricHook(config *eredis_config.Config) redis.Hook {
	return NewRedisHook().SetAfterProcess(
		func(ctx context.Context, cmd redis.Cmder) error {
			duration := GetDuration(ctx)
			err := cmd.Err()
			RedisHandleHistogram.WithLabelValues(config.Name, cmd.Name(), config.AddrString()).Observe(duration.Seconds())
			if err != nil {
				if errors.Is(err, redis.Nil) {
					RedisHandleCounter.WithLabelValues(config.Name, cmd.Name(), config.AddrString(), "NotFound").Inc()
					return err
				}
				RedisHandleCounter.WithLabelValues(config.Name, cmd.Name(), config.AddrString(), "Error").Inc()
				return err
			}

			RedisHandleCounter.WithLabelValues(config.Name, cmd.Name(), config.AddrString(), "OK").Inc()
			return nil
		},
	)
}
