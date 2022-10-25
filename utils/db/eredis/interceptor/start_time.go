package interceptor

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

type startTimeCtxKey struct{}

// StartTimeHook
func StartTimeHook() redis.Hook {
	return NewRedisHook().
		SetBeforeProcess(func(ctx context.Context, cmd redis.Cmder) (context.Context, error) {
			return context.WithValue(ctx, startTimeCtxKey{}, time.Now()), nil
		})
}

func GetStartTime(ctx context.Context) time.Time {
	return ctx.Value(startTimeCtxKey{}).(time.Time)
}

func GetDuration(ctx context.Context) time.Duration {
	startTime, _ := ctx.Value(startTimeCtxKey{}).(time.Time)
	return time.Since(startTime)
}

func GetDurationMilliseconds(ctx context.Context) float64 {
	startTime, _ := ctx.Value(startTimeCtxKey{}).(time.Time)
	return float64(time.Since(startTime).Microseconds()) / 1000

}
