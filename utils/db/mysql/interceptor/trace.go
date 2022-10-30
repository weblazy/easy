package interceptor

import (
	"github.com/weblazy/easy/utils/db/mysql/manager"
	"github.com/weblazy/easy/utils/elog"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"gorm.io/gorm"
)

type TracePlugin struct {
	dsn *manager.DSN
}

func NewTracePlugin(dsn *manager.DSN) *TracePlugin {
	return &TracePlugin{
		dsn: dsn,
	}
}

func (t *TracePlugin) Name() string {
	return "trace"
}

func (t *TracePlugin) Initialize(db *gorm.DB) error {
	var lastErr error
	beforeErrMsg := "TraceStartErr"
	beforeName := "TraceStart"
	beforeFn := t.TraceStart
	err := db.Callback().Query().Before("gorm:query").Register(beforeName, beforeFn)
	if err != nil {
		lastErr = err
		elog.ErrorCtx(db.Statement.Context, beforeErrMsg, zap.Error(err))
	}
	err = db.Callback().Create().Before("gorm:create").Register(beforeName, beforeFn)
	if err != nil {
		lastErr = err
		elog.ErrorCtx(db.Statement.Context, beforeErrMsg, zap.Error(err))
	}
	err = db.Callback().Update().Before("gorm:update").Register(beforeName, beforeFn)
	if err != nil {
		lastErr = err
		elog.ErrorCtx(db.Statement.Context, beforeErrMsg, zap.Error(err))
	}
	err = db.Callback().Delete().Before("gorm:delete").Register(beforeName, beforeFn)
	if err != nil {
		lastErr = err
		elog.ErrorCtx(db.Statement.Context, beforeErrMsg, zap.Error(err))
	}
	err = db.Callback().Row().Before("gorm:row").Register(beforeName, beforeFn)
	if err != nil {
		lastErr = err
		elog.ErrorCtx(db.Statement.Context, beforeErrMsg, zap.Error(err))
	}
	err = db.Callback().Raw().Before("gorm:raw").Register(beforeName, beforeFn)
	if err != nil {
		lastErr = err
		elog.ErrorCtx(db.Statement.Context, beforeErrMsg, zap.Error(err))
	}
	afterErrMsg := "TraceEndErr"
	afterName := "TraceEnd"
	afterFn := t.TraceEnd
	err = db.Callback().Query().After("gorm:query").Register(afterName, afterFn)
	if err != nil {
		lastErr = err
		elog.ErrorCtx(db.Statement.Context, afterErrMsg, zap.Error(err))
	}
	err = db.Callback().Create().After("gorm:create").Register(afterName, afterFn)
	if err != nil {
		lastErr = err
		elog.ErrorCtx(db.Statement.Context, afterErrMsg, zap.Error(err))
	}
	err = db.Callback().Update().After("gorm:update").Register(afterName, afterFn)
	if err != nil {
		lastErr = err
		elog.ErrorCtx(db.Statement.Context, afterErrMsg, zap.Error(err))
	}
	err = db.Callback().Delete().After("gorm:delete").Register(afterName, afterFn)
	if err != nil {
		lastErr = err
		elog.ErrorCtx(db.Statement.Context, afterErrMsg, zap.Error(err))
	}
	err = db.Callback().Row().After("gorm:row").Register(afterName, afterFn)
	if err != nil {
		lastErr = err
		elog.ErrorCtx(db.Statement.Context, afterErrMsg, zap.Error(err))
	}
	err = db.Callback().Raw().After("gorm:raw").Register(afterName, afterFn)
	if err != nil {
		lastErr = err
		elog.ErrorCtx(db.Statement.Context, afterErrMsg, zap.Error(err))
	}
	return lastErr

}

func (t *TracePlugin) TraceStart(db *gorm.DB) {
	tracer := otel.Tracer("")
	if db.Statement.Context != nil {
		_, span := tracer.Start(db.Statement.Context, "GORM", trace.WithSpanKind(trace.SpanKindClient))
		db.InstanceSet("span", span)
		return
	}
}

func (t *TracePlugin) TraceEnd(db *gorm.DB) {
	if db.Statement.Context != nil {
		spanInterface, ok := db.InstanceGet("span")
		if !ok || spanInterface == nil {
			return
		}
		span := spanInterface.(trace.Span)
		defer span.End()
		span.SetAttributes(
			attribute.String("sql.inner", t.dsn.DBName),
			attribute.String("sql.addr", t.dsn.Addr),
			attribute.String("peer.service", "mysql"),
			attribute.String("db.instance", t.dsn.DBName),
			attribute.String("peer.address", t.dsn.Addr),
			attribute.String("peer.statement", logSQL(db.Statement.SQL.String(), db.Statement.Vars, false)),
		)
		return
	}

}
