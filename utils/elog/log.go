package elog

import (
	"context"
	"sync"

	"github.com/weblazy/easy/utils/elog/logx"
)

var (
	Logger sync.Map
)

type LoggerNameCtxKey struct{}

// 默认加入zap组件
func init() {
	Logger.Store(Ezap, DefaultLogger)
}

// SetLogger 设置日志打印实例,选择输出到文件,终端,阿里云日志等
func SetLogger(name string, logger logx.GLog) {
	Logger.Store(name, logger)
}

// DelLogger 删除日志插件
func DelLogger(name string) {
	Logger.Delete(name)
}

// GetLoggerFromCtx
func GetLoggerFromCtx(ctx context.Context) logx.GLog {
	loggerName, ok := ctx.Value(LoggerNameCtxKey{}).(string)
	if ok {
		logger, ok := Logger.Load(loggerName)
		if !ok {
			// 指定了logger,但是没有找到
			return nil
		}
		return logger.(logx.GLog)
	}
	// 没有指定logger使用默认全局logger
	return DefaultLogger
}

func SetLogerName(ctx context.Context, name string) context.Context {
	return context.WithValue(ctx, LoggerNameCtxKey{}, name)
}
