package mysql

import (
	"context"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/weblazy/easy/utils/glog"
	"go.uber.org/zap"
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
	for {
		time.Sleep(time.Second * 10)
		_Gorm.gormMaps.Range(func(key, value interface{}) bool {
			name := key.(string)
			db := value.(*gorm.DB)
			sqlDB, err := db.DB()
			if err != nil {
				glog.ErrorCtx(context.Background(), "monitor db error", zap.Error(err))
				return false
			}

			stats := sqlDB.Stats()

			DBStatsGauge.WithLabelValues(name, "idle").Set(float64(stats.Idle))
			DBStatsGauge.WithLabelValues(name, "inuse").Set(float64(stats.InUse))
			DBStatsGauge.WithLabelValues(name, "wait").Set(float64(stats.WaitCount))
			DBStatsGauge.WithLabelValues(name, "conns").Set(float64(stats.OpenConnections))
			DBStatsGauge.WithLabelValues(name, "max_open_conns").Set(float64(stats.MaxOpenConnections))
			DBStatsGauge.WithLabelValues(name, "max_idle_closed").Set(float64(stats.MaxIdleClosed))
			DBStatsGauge.WithLabelValues(name, "max_lifetime_closed").Set(float64(stats.MaxLifetimeClosed))

			return true
		})
	}
}
