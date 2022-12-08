package interceptor

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

type ctxStartTimeKey struct{}

// StartTimeHook
func StartTimeHook() redis.Hook {
	return NewRedisHook().
		SetBeforeProcess(func(ctx context.Context, cmd redis.Cmder) (context.Context, error) {
			return context.WithValue(ctx, ctxStartTimeKey{}, time.Now()), nil
		})
}

func GetStartTime(ctx context.Context) time.Time {
	return ctx.Value(ctxStartTimeKey{}).(time.Time)
}

func GetDuration(ctx context.Context) time.Duration {
	startTime, _ := ctx.Value(ctxStartTimeKey{}).(time.Time)
	return time.Since(startTime)
}

func GetDurationMilliseconds(ctx context.Context) float64 {
	startTime, _ := ctx.Value(ctxStartTimeKey{}).(time.Time)
	return float64(time.Since(startTime).Microseconds()) / 1000

}
