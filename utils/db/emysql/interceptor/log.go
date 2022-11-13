package interceptor

import (
	"context"
	"errors"
	"time"

	"github.com/weblazy/easy/utils/db/emysql/manager"
	"github.com/weblazy/easy/utils/db/emysql/mysql_config"
	"github.com/weblazy/easy/utils/elog"
	"github.com/weblazy/easy/utils/etrace"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type LogPlugin struct {
	config *mysql_config.Config
	dsn    *manager.DSN
}

func NewLogPlugin(config *mysql_config.Config, dsn *manager.DSN) *LogPlugin {
	return &LogPlugin{
		config: config,
		dsn:    dsn,
	}
}

func (e *LogPlugin) Name() string {
	return "log"
}

func (e *LogPlugin) Initialize(db *gorm.DB) error {
	var lastErr error
	afterErrMsg := "LogEndErr"
	afterName := "LogEnd"
	err := db.Callback().Query().After("gorm:query").Register(afterName, e.LogEnd("gorm:query"))
	if err != nil {
		lastErr = err
		elog.ErrorCtx(db.Statement.Context, afterErrMsg, zap.Error(err))
	}
	err = db.Callback().Create().After("gorm:create").Register(afterName, e.LogEnd("gorm:create"))
	if err != nil {
		lastErr = err
		elog.ErrorCtx(db.Statement.Context, afterErrMsg, zap.Error(err))
	}
	err = db.Callback().Update().After("gorm:update").Register(afterName, e.LogEnd("gorm:update"))
	if err != nil {
		lastErr = err
		elog.ErrorCtx(db.Statement.Context, afterErrMsg, zap.Error(err))
	}
	err = db.Callback().Delete().After("gorm:delete").Register(afterName, e.LogEnd("gorm:delete"))
	if err != nil {
		lastErr = err
		elog.ErrorCtx(db.Statement.Context, afterErrMsg, zap.Error(err))
	}
	err = db.Callback().Row().After("gorm:row").Register(afterName, e.LogEnd("gorm:row"))
	if err != nil {
		lastErr = err
		elog.ErrorCtx(db.Statement.Context, afterErrMsg, zap.Error(err))
	}
	err = db.Callback().Raw().After("gorm:raw").Register(afterName, e.LogEnd("gorm:raw"))
	if err != nil {
		lastErr = err
		elog.ErrorCtx(db.Statement.Context, afterErrMsg, zap.Error(err))
	}
	return lastErr

}

func (e *LogPlugin) LogEnd(method string) func(db *gorm.DB) {
	return func(db *gorm.DB) {
		ctx := db.Statement.Context
		duration := GetDuration(ctx)
		//自定义key todo
		var loggerKeys []string
		var fields = make([]zap.Field, 0, 15+len(loggerKeys))
		fields = append(fields,
			elog.FieldMethod(method),
			elog.FieldName(e.dsn.DBName+"."+db.Statement.Table), elog.FieldCost(duration))

		if e.config.EnableAccessInterceptorReq {
			// todo: EnableDetailSQL 参数是否只在错误时生效
			fields = append(fields, zap.String(elog.KeyReq, logSQL(db.Statement.SQL.String(), db.Statement.Vars, e.config.EnableDetailSQL)))
		}

		if e.config.EnableAccessInterceptorRes {
			fields = append(fields, elog.FieldResp(db.Statement.Dest))
		}

		// 开启了链路，那么就记录链路id
		if e.config.EnableTraceInterceptor {
			fields = append(fields, elog.FieldTrace(etrace.ExtractTraceID(db.Statement.Context)))
		}

		// 支持自定义log
		for _, key := range loggerKeys {
			if value := getContextValue(db.Statement.Context, key); value != "" {
				fields = append(fields, zap.String(key, value))
			}
		}

		// 如果有慢日志，就记录
		var isSlow bool
		if e.config.SlowLogThreshold > time.Duration(0) && e.config.SlowLogThreshold < duration {
			isSlow = true
		}
		fields = append(fields, elog.FieldSlow(isSlow))
		// 如果有错误，记录错误信息
		if db.Error != nil {
			fields = append(fields, elog.FieldEvent("error"), elog.FieldError(db.Error))
			if errors.Is(db.Error, ErrRecordNotFound) {
				if e.config.EnableRecordNotFoundLog {
					elog.WarnCtx(db.Statement.Context, mysql_config.PkgName, fields...)
				}
				return
			}
			elog.ErrorCtx(db.Statement.Context, mysql_config.PkgName, fields...)
			return
		}

		if isSlow {
			elog.WarnCtx(db.Statement.Context, mysql_config.PkgName, fields...)
			return
		}

		// 开启了记录日志信息，那么就记录access
		// event normal和error，代表全部access的请求数
		if e.config.EnableAccessInterceptor {
			elog.InfoCtx(db.Statement.Context, mysql_config.PkgName, fields...)
		}

	}
}

func getContextValue(c context.Context, key string) string {
	v, _ := c.Value(key).(string)
	return v
}
