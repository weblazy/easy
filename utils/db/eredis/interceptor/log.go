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
			duration := GetDuration(ctx)
			fields = append(fields, elog.FieldName(config.Name),
				elog.FieldMethod(cmd.Name()),
				elog.FieldDuration(duration))

			if config.EnableLogReq {
				fields = append(fields, elog.FieldReq(cmd.Args()))
			}
			if config.EnableLogRes && err == nil {
				fields = append(fields, elog.FieldResp(response(cmd)))
			}

			// 开启了链路，那么就记录链路id
			if config.EnableTraceInterceptor && etrace.IsGlobalTracerRegistered() {
				fields = append(fields, elog.FieldTrace(etrace.ExtractTraceID(ctx)))
			}
			var isSlow bool
			if config.SlowLogThreshold > time.Duration(0) && duration > config.SlowLogThreshold {
				isSlow = true
			}
			fields = append(fields, elog.FieldSlow(isSlow))
			if err != nil {
				fields = append(fields, elog.FieldError(err))
				if errors.Is(err, redis.Nil) {
					elog.WarnCtx(ctx, eredis_config.PkgName, fields...)
					return err
				}
				elog.ErrorCtx(ctx, eredis_config.PkgName, fields...)
				return err
			}
			if isSlow {
				elog.WarnCtx(ctx, eredis_config.PkgName, fields...)
				return nil
			}
			if config.EnableLogAccess {
				elog.InfoCtx(ctx, eredis_config.PkgName, fields...)
			}
			return nil
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
