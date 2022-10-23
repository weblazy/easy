package interceptor

import (
	"github.com/weblazy/easy/utils/db/mysql/manager"
	"go.opentelemetry.io/otel/trace"

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
	return "metric"
}

func (t *TracePlugin) Initialize(db *gorm.DB) error {
	return db.Callback().Query().After("gorm:query").Register("explain", t.TraceStart)
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
