package ezap

import (
	"context"
	"os"

	"github.com/weblazy/easy/utils/elog/logx"

	"go.uber.org/zap"
)

// Ezap 将文件输出到终端或者文件
type Ezap struct {
	logx.GLog
	Logger  *zap.Logger
	Logfile *os.File
	Config  *Config
}

func (e *Ezap) DebugCtx(ctx context.Context, msg string, fields ...zap.Field) {
	e.Logger.Debug(msg, fields...)
}

func (e *Ezap) InfoCtx(ctx context.Context, msg string, fields ...zap.Field) {
	e.Logger.Info(msg, fields...)
}

func (e *Ezap) WarnCtx(ctx context.Context, msg string, fields ...zap.Field) {
	e.Logger.Warn(msg, fields...)
}

func (e *Ezap) ErrorCtx(ctx context.Context, msg string, fields ...zap.Field) {
	e.Logger.Error(msg, fields...)
}
