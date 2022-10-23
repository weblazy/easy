package interceptor

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/weblazy/easy/utils/db/eredis/eredis_config"
	"github.com/weblazy/easy/utils/elog"
	"github.com/weblazy/easy/utils/etrace"
	"go.uber.org/zap"
)

func LogHook(config *eredis_config.Config) redis.Hook {
	return NewRedisHook().SetAfterProcess(
		func(ctx context.Context, cmd redis.Cmder) error {
			var fields = make([]zap.Field, 0, 15)
			var err = cmd.Err()
			duration := time.Since(ctx.Value(startTimeCtxKey{}).(time.Time))
			fields = append(fields, zap.String("name", config.Name),
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
				fields = append(fields, zap.Bool("slow", true))
			}

			// error metric
			if err != nil {
				fields = append(fields, zap.String("event", "error"), zap.Error(err))
				if errors.Is(err, redis.Nil) {
					elog.WarnCtx(ctx, "eredis", fields...)
					return err
				}
				elog.ErrorCtx(ctx, "eredis", fields...)
				return err
			}

			if config.EnableAccessInterceptor {
				fields = append(fields, zap.String("event", "normal"))
				elog.InfoCtx(ctx, "eredis", fields...)
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
