package interceptor

import (
	"context"
	"time"

	"github.com/weblazy/easy/utils/elog"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type StartTimePlugin struct{}

type startTimeCtxkey struct{}

func NewStartTimePlugin() *StartTimePlugin {
	return &StartTimePlugin{}
}

func (e *StartTimePlugin) Name() string {
	return "start_time"
}

func (e *StartTimePlugin) Initialize(db *gorm.DB) error {
	var lastErr error
	err := db.Callback().Query().Before("gorm:query").Register("SetStartTime", SetStartTime)
	if err != nil {
		lastErr = err
		elog.ErrorCtx(db.Statement.Context, "SetStartTimeErr", zap.Error(err))
	}
	err = db.Callback().Create().Before("gorm:create").Register("SetStartTime", SetStartTime)
	if err != nil {
		lastErr = err
		elog.ErrorCtx(db.Statement.Context, "SetStartTimeErr", zap.Error(err))
	}
	err = db.Callback().Update().Before("gorm:update").Register("SetStartTime", SetStartTime)
	if err != nil {
		lastErr = err
		elog.ErrorCtx(db.Statement.Context, "SetStartTimeErr", zap.Error(err))
	}
	err = db.Callback().Delete().Before("gorm:delete").Register("SetStartTime", SetStartTime)
	if err != nil {
		lastErr = err
		elog.ErrorCtx(db.Statement.Context, "SetStartTimeErr", zap.Error(err))
	}
	err = db.Callback().Query().Before("gorm:query").Register("SetStartTime", SetStartTime)
	if err != nil {
		lastErr = err
		elog.ErrorCtx(db.Statement.Context, "SetStartTimeErr", zap.Error(err))
	}
	err = db.Callback().Row().Before("gorm:row").Register("SetStartTime", SetStartTime)
	if err != nil {
		lastErr = err
		elog.ErrorCtx(db.Statement.Context, "SetStartTimeErr", zap.Error(err))
	}
	err = db.Callback().Raw().Before("gorm:raw").Register("SetStartTime", SetStartTime)
	if err != nil {
		lastErr = err
		elog.ErrorCtx(db.Statement.Context, "SetStartTimeErr", zap.Error(err))
	}
	return lastErr
}

func SetStartTime(db *gorm.DB) {
	startTime := time.Now()
	db.Statement.Context = context.WithValue(db.Statement.Context, startTimeCtxkey{}, startTime)
	return
}

func GetStartTime(db *gorm.DB) time.Time {
	return db.Statement.Context.Value(startTimeCtxkey{}).(time.Time)
}

func GetDuration(ctx context.Context) time.Duration {
	startTime, _ := ctx.Value(startTimeCtxkey{}).(time.Time)
	return time.Since(startTime)
}

func GetDurationMilliseconds(ctx context.Context) float64 {
	startTime, _ := ctx.Value(startTimeCtxkey{}).(time.Time)
	return float64(time.Since(startTime).Microseconds()) / 1000

}
