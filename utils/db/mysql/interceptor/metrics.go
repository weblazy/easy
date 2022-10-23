package interceptor

import (
	"errors"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/weblazy/easy/utils/db/mysql/manager"
	"github.com/weblazy/easy/utils/db/mysql/mysql_config"
	"gorm.io/gorm"
)

var (
	DBHandleHistogram = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name: "db_handle_seconds",
	}, []string{"type", "name", "method", "peer"})

	DBHandleCounter = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "db_handle_total",
	}, []string{"type", "name", "method", "peer", "code"})

	DBStatsGauge = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "db_stats",
	}, []string{"name", "type"})
)

func init() {
	prometheus.MustRegister(DBHandleHistogram)
	prometheus.MustRegister(DBHandleCounter)
	prometheus.MustRegister(DBStatsGauge)
	//nolint:gochecknoinits
	go monitor()
}

func monitor() {
	// for {
	// 	time.Sleep(time.Second * 10)
	// 	_Gorm.gormMaps.Range(func(key, value interface{}) bool {
	// 		name := key.(string)
	// 		db := value.(*gorm.DB)
	// 		sqlDB, err := db.DB()
	// 		if err != nil {
	// 			elog.ErrorCtx(context.Background(), "monitor db error", zap.Error(err))
	// 			return false
	// 		}

	// 		stats := sqlDB.Stats()

	// 		DBStatsGauge.WithLabelValues(name, "idle").Set(float64(stats.Idle))
	// 		DBStatsGauge.WithLabelValues(name, "inuse").Set(float64(stats.InUse))
	// 		DBStatsGauge.WithLabelValues(name, "wait").Set(float64(stats.WaitCount))
	// 		DBStatsGauge.WithLabelValues(name, "conns").Set(float64(stats.OpenConnections))
	// 		DBStatsGauge.WithLabelValues(name, "max_open_conns").Set(float64(stats.MaxOpenConnections))
	// 		DBStatsGauge.WithLabelValues(name, "max_idle_closed").Set(float64(stats.MaxIdleClosed))
	// 		DBStatsGauge.WithLabelValues(name, "max_lifetime_closed").Set(float64(stats.MaxLifetimeClosed))

	// 		return true
	// 	})
	// }
}

type MetricPlugin struct {
	dsn    *manager.DSN
	config *mysql_config.Config
}

func NewMetricPlugin() *ExplainPlugin {
	return &ExplainPlugin{}
}

func (e *MetricPlugin) Name() string {
	return "metric"
}

func (e *MetricPlugin) Initialize(db *gorm.DB) error {
	return db.Callback().Query().After("gorm:query").Register("MetricEnd", e.MetricEnd("gorm:query"))
}

func (e *MetricPlugin) MetricEnd(method string) func(db *gorm.DB) {
	return func(db *gorm.DB) {
		ctx := db.Statement.Context
		duration := GetDuration(ctx)

		// 记录监控耗时
		DBHandleHistogram.WithLabelValues(TypeGorm, e.config.Name, e.dsn.DBName+"."+db.Statement.Table, e.dsn.Addr).Observe(duration.Seconds())

		// 如果有错误，记录错误信息
		if db.Error != nil {
			if errors.Is(db.Error, ErrRecordNotFound) {
				DBHandleCounter.WithLabelValues(TypeGorm, e.config.Name, e.dsn.DBName+"."+db.Statement.Table, e.dsn.Addr, "Empty").Inc()
				return
			}
			DBHandleCounter.WithLabelValues(TypeGorm, e.config.Name, e.dsn.DBName+"."+db.Statement.Table, e.dsn.Addr, "Error").Inc()
			return
		}

		DBHandleCounter.WithLabelValues(TypeGorm, e.config.Name, e.dsn.DBName+"."+db.Statement.Table, e.dsn.Addr, "OK").Inc()
	}
}
