package interceptor

import (
	"context"
	"fmt"
	"time"

	"errors"

	"github.com/go-redis/redis/v8"

	"github.com/weblazy/easy/utils/db/eredis/eredis_config"
	"github.com/weblazy/easy/utils/elog"
	"github.com/weblazy/easy/utils/etrace"
	"go.uber.org/zap"
)

// https://stackoverflow.com/questions/40891345/fix-should-not-use-basic-type-string-as-key-in-context-withvalue-golint
// https://blog.golang.org/context#TOC_3.2.
// https://golang.org/pkg/context/#WithValue ，这边文章说明了用struct，可以避免分配
type fredis2ContextKeyType struct{}

var ctxBegKey = fredis2ContextKeyType{}

type Interceptor struct {
	beforeProcess         func(ctx context.Context, cmd redis.Cmder) (context.Context, error)
	afterProcess          func(ctx context.Context, cmd redis.Cmder) error
	beforeProcessPipeline func(ctx context.Context, cmds []redis.Cmder) (context.Context, error)
	afterProcessPipeline  func(ctx context.Context, cmds []redis.Cmder) error
}

func (i *Interceptor) BeforeProcess(ctx context.Context, cmd redis.Cmder) (context.Context, error) {
	return i.beforeProcess(ctx, cmd)
}

func (i *Interceptor) AfterProcess(ctx context.Context, cmd redis.Cmder) error {
	return i.afterProcess(ctx, cmd)
}

func (i *Interceptor) BeforeProcessPipeline(ctx context.Context, cmds []redis.Cmder) (context.Context, error) {
	return i.beforeProcessPipeline(ctx, cmds)
}

func (i *Interceptor) AfterProcessPipeline(ctx context.Context, cmds []redis.Cmder) error {
	return i.afterProcessPipeline(ctx, cmds)
}

func newInterceptor(compName string, config *eredis_config.Config) *Interceptor {
	return &Interceptor{
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

func (i *Interceptor) setBeforeProcess(p func(ctx context.Context, cmd redis.Cmder) (context.Context, error)) *Interceptor {
	i.beforeProcess = p
	return i
}

func (i *Interceptor) setAfterProcess(p func(ctx context.Context, cmd redis.Cmder) error) *Interceptor {
	i.afterProcess = p
	return i
}

func (i *Interceptor) setBeforeProcessPipeline(p func(ctx context.Context, cmds []redis.Cmder) (context.Context, error)) *Interceptor { //nolint
	i.beforeProcessPipeline = p
	return i
}

func (i *Interceptor) setAfterProcessPipeline(p func(ctx context.Context, cmds []redis.Cmder) error) *Interceptor { //nolint
	i.afterProcessPipeline = p
	return i
}

func FixedInterceptor(compName string, config *eredis_config.Config) *Interceptor {
	return newInterceptor(compName, config).
		setBeforeProcess(func(ctx context.Context, cmd redis.Cmder) (context.Context, error) {
			return context.WithValue(ctx, ctxBegKey, time.Now()), nil
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

func DebugInterceptor(compName string, config *eredis_config.Config) *Interceptor {
	return newInterceptor(compName, config).setAfterProcess(
		func(ctx context.Context, cmd redis.Cmder) error {
			duration := time.Since(ctx.Value(ctxBegKey).(time.Time))
			err := cmd.Err()
			if err != nil {
				elog.ErrorCtx(ctx, "eredis.response", elog.MakeReqResError(1, compName, config.AddrString(), duration, fmt.Sprintf("%v", cmd.Args()), err.Error()))

			} else {
				elog.InfoCtx(ctx, "eredis.response", elog.MakeReqResInfo(1, compName, config.AddrString(), duration, fmt.Sprintf("%v", cmd.Args()), response(cmd)))
			}
			return err
		},
	)
}

func MetricInterceptor(compName string, config *eredis_config.Config) *Interceptor {
	return newInterceptor(compName, config).setAfterProcess(
		func(ctx context.Context, cmd redis.Cmder) error {
			duration := time.Since(ctx.Value(ctxBegKey).(time.Time))
			err := cmd.Err()
			RedisHandleHistogram.WithLabelValues(compName, cmd.Name(), config.AddrString()).Observe(duration.Seconds())
			if err != nil {
				if errors.Is(err, redis.Nil) {
					RedisHandleCounter.WithLabelValues(compName, cmd.Name(), config.AddrString(), "Empty").Inc()
					return err
				}
				RedisHandleCounter.WithLabelValues(compName, cmd.Name(), config.AddrString(), "Error").Inc()
				return err
			}

			RedisHandleCounter.WithLabelValues(compName, cmd.Name(), config.AddrString(), "OK").Inc()
			return nil
		},
	)
}

func AccessInterceptor(compName string, config *eredis_config.Config) *Interceptor {
	return newInterceptor(compName, config).setAfterProcess(
		func(ctx context.Context, cmd redis.Cmder) error {
			var fields = make([]zap.Field, 0, 15)
			var err = cmd.Err()
			duration := time.Since(ctx.Value(ctxBegKey).(time.Time))
			fields = append(fields, zap.String("name", compName),
				zap.String("method", cmd.Name()),
				zap.Int64("duration", duration.Microseconds()/1000))

			if config.EnableAccessInterceptorReq {
				fields = append(fields, zap.Any("req", cmd.Args()))
			}
			if config.EnableAccessInterceptorRes && err == nil {
				fields = append(fields, zap.String("res", response(cmd)))
			}

			// 开启了链路，那么就记录链路id
			if config.EnableTraceInterceptor && etrace.IsGlobalTracerRegistered() {
				fields = append(fields, zap.String("trace_id", etrace.ExtractTraceID(ctx)))
			}

			if config.SlowLogThreshold > time.Duration(0) && duration > config.SlowLogThreshold {
				elog.InfoCtx(ctx, "slow", fields...)
			}

			// error metric
			if err != nil {
				fields = append(fields, zap.String("event", "error"), zap.Error(err))
				if errors.Is(err, redis.Nil) {
					elog.WarnCtx(ctx, "access", fields...)
					return err
				}
				elog.ErrorCtx(ctx, "access", fields...)
				return err
			}

			if config.EnableAccessInterceptor {
				fields = append(fields, zap.String("event", "normal"))
				elog.InfoCtx(ctx, "access", fields...)
			}
			return err
		},
	)
}

func response(cmd redis.Cmder) string {
	switch t := cmd.(type) {
	case *redis.Cmd:
		return fmt.Sprintf("%v", t.Val())
	case *redis.StringCmd:
		return t.Val()
	case *redis.StatusCmd:
		return t.Val()
	case *redis.IntCmd:
		return fmt.Sprintf("%v", t.Val())
	case *redis.DurationCmd:
		return t.Val().String()
	case *redis.BoolCmd:
		return fmt.Sprintf("%v", t.Val())
	case *redis.CommandsInfoCmd:
		return fmt.Sprintf("%v", t.Val())
	case *redis.StringSliceCmd:
		return fmt.Sprintf("%v", t.Val())
	default:
		return ""
	}
}
