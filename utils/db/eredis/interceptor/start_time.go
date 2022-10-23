package interceptor

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

func StartTimeHook() redis.Hook {
	return NewRedisHook().
		SetBeforeProcess(func(ctx context.Context, cmd redis.Cmder) (context.Context, error) {
			return context.WithValue(ctx, startTimeCtxKey{}, time.Now()), nil
		})
	// 这里会改写错误
	//setAfterProcess(func(ctx context.Context, cmd redis.Cmder) error {
	//	var err = cmd.Err()
	//	// go-redis script的error做了prefix处理
	//	// https://github.com/go-redis/redis/blob/master/script.go#L61
	//	if err != nil && !strings.HasPrefix(err.Error(), "NOSCRIPT ") {
	//		err = fmt.Errorf("eredis exec command %s fail, %w", cmd.Name(), err)
	//	}
	//	return err
	//})
}
