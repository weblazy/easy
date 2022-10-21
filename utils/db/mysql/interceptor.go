package mysql

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/weblazy/easy/utils/elog"
	"github.com/weblazy/easy/utils/etrace"
	"go.uber.org/zap"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	"github.com/weblazy/easy/utils/db/mysql/manager"

	"gorm.io/gorm"
)

const (
	TypeGorm = "gorm"
)

// Handler ...
type Handler func(*gorm.DB)

// Processor ...
type Processor interface {
	Get(name string) func(*gorm.DB)
	Replace(name string, handler func(*gorm.DB)) error
}

// Interceptor ...
type Interceptor func(string, *manager.DSN, string, *Config) func(next Handler) Handler

// 确保在生产不要开 debug
func debugInterceptor(compName string, dsn *manager.DSN, op string, options *Config) func(Handler) Handler {
	return func(next Handler) Handler {
		return func(db *gorm.DB) {
			beg := time.Now()
			next(db)
			duration := time.Since(beg)
			if db.Error != nil {
				elog.ErrorCtx(db.Statement.Context, "fgorm.response", elog.MakeReqResError(1, compName, dsn.Addr+"/"+dsn.DBName, duration, logSQL(db.Statement.SQL.String(), db.Statement.Vars, true), db.Error.Error()))
			} else {
				elog.InfoCtx(db.Statement.Context, "fgorm.response", elog.MakeReqResInfo(1, compName, dsn.Addr+"/"+dsn.DBName, duration, logSQL(db.Statement.SQL.String(), db.Statement.Vars, true), fmt.Sprintf("%v", db.Statement.Dest)))
			}
		}
	}
}

func metricInterceptor(compName string, dsn *manager.DSN, op string, config *Config) func(Handler) Handler {
	return func(next Handler) Handler {
		return func(db *gorm.DB) {
			beg := time.Now()
			next(db)
			duration := time.Since(beg)

			var loggerKeys []string

			var fields = make([]zap.Field, 0, 15+len(loggerKeys))
			fields = append(fields,
				elog.FieldMethod(op),
				elog.FieldName(dsn.DBName+"."+db.Statement.Table), elog.FieldCost(duration))

			if config.EnableAccessInterceptorReq {
				// todo: EnableDetailSQL 参数是否只在错误时生效
				fields = append(fields, zap.String(elog.KeyReq, logSQL(db.Statement.SQL.String(), db.Statement.Vars, config.EnableDetailSQL)))
			}

			if config.EnableAccessInterceptorRes {
				fields = append(fields, elog.FieldResp(db.Statement.Dest))
			}

			// 开启了链路，那么就记录链路id
			if config.EnableTraceInterceptor {
				fields = append(fields, elog.FieldTrace(etrace.ExtractTraceID(db.Statement.Context)))
			}

			// 支持自定义log
			for _, key := range loggerKeys {
				if value := getContextValue(db.Statement.Context, key); value != "" {
					fields = append(fields, zap.String(key, value))
				}
			}

			// 记录监控耗时
			DBHandleHistogram.WithLabelValues(TypeGorm, compName, dsn.DBName+"."+db.Statement.Table, dsn.Addr).Observe(duration.Seconds())

			// 如果有慢日志，就记录
			if config.SlowLogThreshold > time.Duration(0) && config.SlowLogThreshold < duration {
				elog.WarnCtx(db.Statement.Context, "slow", fields...)
			}

			// 如果有错误，记录错误信息
			if db.Error != nil {
				fields = append(fields, elog.FieldEvent("error"), elog.FieldError(db.Error))
				if errors.Is(db.Error, ErrRecordNotFound) {
					if config.EnableRecordNotFoundLog {
						elog.WarnCtx(db.Statement.Context, "access", fields...)
					}
					DBHandleCounter.WithLabelValues(TypeGorm, compName, dsn.DBName+"."+db.Statement.Table, dsn.Addr, "Empty").Inc()
					return
				}
				elog.ErrorCtx(db.Statement.Context, "access", fields...)
				DBHandleCounter.WithLabelValues(TypeGorm, compName, dsn.DBName+"."+db.Statement.Table, dsn.Addr, "Error").Inc()
				return
			}

			DBHandleCounter.WithLabelValues(TypeGorm, compName, dsn.DBName+"."+db.Statement.Table, dsn.Addr, "OK").Inc()
			// 开启了记录日志信息，那么就记录access
			// event normal和error，代表全部access的请求数
			if config.EnableAccessInterceptor {
				fields = append(fields,
					elog.FieldEvent("normal"),
				)
				elog.InfoCtx(db.Statement.Context, "access", fields...)
			}
		}
	}
}

func logSQL(sql string, args []interface{}, containArgs bool) string {
	if containArgs {
		return bindSQL(sql, args)
	}
	return sql
}

func traceInterceptor(compName string, dsn *manager.DSN, op string, options *Config) func(Handler) Handler {
	tracer := otel.Tracer("")

	return func(next Handler) Handler {
		return func(db *gorm.DB) {
			if db.Statement.Context != nil {
				_, span := tracer.Start(db.Statement.Context, "GORM", trace.WithSpanKind(trace.SpanKindClient))
				defer span.End()
				// 延迟执行 scope.CombinedConditionSql() 避免sqlVar被重复追加
				next(db)

				span.SetAttributes(
					attribute.String("sql.inner", dsn.DBName),
					attribute.String("sql.addr", dsn.Addr),
					attribute.String("peer.service", "mysql"),
					attribute.String("db.instance", dsn.DBName),
					attribute.String("peer.address", dsn.Addr),
					attribute.String("peer.statement", logSQL(db.Statement.SQL.String(), db.Statement.Vars, false)),
				)
				return
			}

			next(db)
		}
	}
}

func getContextValue(c context.Context, key string) string {
	if key == "" {
		return ""
	}
	return c.Value(key).(string)
}
