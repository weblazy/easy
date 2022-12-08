package interceptor

import (
	"context"

	"github.com/go-redis/redis/v8"
)

// https://stackoverflow.com/questions/40891345/fix-should-not-use-basic-type-string-as-key-in-context-withvalue-golint
// https://blog.golang.org/context#TOC_3.2.
// https://golang.org/pkg/context/#WithValue ，这边文章说明了用struct，可以避免分配

type RedisHook struct {
	redis.Hook
	beforeProcess         func(ctx context.Context, cmd redis.Cmder) (context.Context, error)
	afterProcess          func(ctx context.Context, cmd redis.Cmder) error
	beforeProcessPipeline func(ctx context.Context, cmds []redis.Cmder) (context.Context, error)
	afterProcessPipeline  func(ctx context.Context, cmds []redis.Cmder) error
}

func (i *RedisHook) BeforeProcess(ctx context.Context, cmd redis.Cmder) (context.Context, error) {
	return i.beforeProcess(ctx, cmd)
}

func (i *RedisHook) AfterProcess(ctx context.Context, cmd redis.Cmder) error {
	return i.afterProcess(ctx, cmd)
}

func (i *RedisHook) BeforeProcessPipeline(ctx context.Context, cmds []redis.Cmder) (context.Context, error) {
	return i.beforeProcessPipeline(ctx, cmds)
}

func (i *RedisHook) AfterProcessPipeline(ctx context.Context, cmds []redis.Cmder) error {
	return i.afterProcessPipeline(ctx, cmds)
}

func NewRedisHook() *RedisHook {
	return &RedisHook{
		beforeProcess: func(ctx context.Context, cmd redis.Cmder) (context.Context, error) {
			return ctx, nil
		},
		afterProcess: func(ctx context.Context, cmd redis.Cmder) error {
			return nil
		},
		beforeProcessPipeline: func(ctx context.Context, cmds []redis.Cmder) (context.Context, error) {
			return ctx, nil
		},
		afterProcessPipeline: func(ctx context.Context, cmds []redis.Cmder) error {
			return nil
		},
	}
}

func (i *RedisHook) SetBeforeProcess(p func(ctx context.Context, cmd redis.Cmder) (context.Context, error)) *RedisHook {
	i.beforeProcess = p
	return i
}

func (i *RedisHook) SetAfterProcess(p func(ctx context.Context, cmd redis.Cmder) error) *RedisHook {
	i.afterProcess = p
	return i
}

func (i *RedisHook) SetBeforeProcessPipeline(p func(ctx context.Context, cmds []redis.Cmder) (context.Context, error)) *RedisHook { //nolint
	i.beforeProcessPipeline = p
	return i
}

func (i *RedisHook) SetAfterProcessPipeline(p func(ctx context.Context, cmds []redis.Cmder) error) *RedisHook { //nolint
	i.afterProcessPipeline = p
	return i
}
