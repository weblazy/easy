package interceptor

import (
	"context"
	"time"

	"github.com/weblazy/easy/utils/elog"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type StartTimePlugin struct{}

type ctxStartTimeKey struct{}

func NewStartTimePlugin() *StartTimePlugin {
	return &StartTimePlugin{}
}

func (e *StartTimePlugin) Name() string {
	return "start_time"
}

func (e *StartTimePlugin) Initialize(db *gorm.DB) error {
	var lastErr error
	beforeErrMsg := "SetStartTimeErr"
	beforeName := "SetStartTime"
	beforeFn := SetStartTime
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
	return lastErr
}

func SetStartTime(db *gorm.DB) {
	startTime := time.Now()
	db.Statement.Context = context.WithValue(db.Statement.Context, ctxStartTimeKey{}, startTime)
}

func GetStartTime(db *gorm.DB) time.Time {
	return db.Statement.Context.Value(ctxStartTimeKey{}).(time.Time)
}

func GetDuration(ctx context.Context) time.Duration {
	startTime, _ := ctx.Value(ctxStartTimeKey{}).(time.Time)
	return time.Since(startTime)
}

func GetDurationMilliseconds(ctx context.Context) float64 {
	startTime, _ := ctx.Value(ctxStartTimeKey{}).(time.Time)
	return float64(time.Since(startTime).Microseconds()) / 1000

}
