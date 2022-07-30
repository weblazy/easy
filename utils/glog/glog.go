package glog

import (
	uzap "go.uber.org/zap"
	"golang.org/x/net/context"

	"github.com/weblazy/easy/utils/glog/logx"
)

const CtxKey = "logConf"

type LogConf struct {
	Name   string
	Labels []uzap.Field
}

var defaultLogConf = &LogConf{}

func GetContextLog(ctx context.Context) *LogConf {
	if v, ok := ctx.Value(CtxKey).(*LogConf); !ok {
		return v
	} else {
		return defaultLogConf
	}
}

func GetLabels(ctx context.Context, fields ...uzap.Field) []uzap.Field {
	logConf := GetContextLog(ctx)
	llen := len(logConf.Labels)
	flen := len(fields)
	newFields := make([]uzap.Field, llen+flen)
	copy(newFields, logConf.Labels)
	copy(newFields[llen:], newFields)
	return newFields
}

func InfoCtx(ctx context.Context, msg string, fields ...uzap.Field) {
	fields = GetLabels(ctx, fields...)
	Logger.Range(func(k, v interface{}) bool {
		v.(logx.GLog).Info(msg, fields...)
		return true
	})
}

func InfoCtxF(ctx context.Context, format string, args ...interface{}) {
	Logger.Range(func(k, v interface{}) bool {
		v.(logx.GLog).InfoF(format, args...)
		return true
	})
}

func DebugCtx(ctx context.Context, msg string, fields ...uzap.Field) {
	fields = GetLabels(ctx, fields...)
	Logger.Range(func(k, v interface{}) bool {
		v.(logx.GLog).Debug(msg, fields...)
		return true
	})
}

func DebugCtxF(ctx context.Context, format string, args ...interface{}) {
	Logger.Range(func(k, v interface{}) bool {
		v.(logx.GLog).DebugF(format, args...)
		return true
	})
}

func WarnCtx(ctx context.Context, msg string, fields ...uzap.Field) {
	fields = GetLabels(ctx, fields...)
	Logger.Range(func(k, v interface{}) bool {
		v.(logx.GLog).Warn(msg, fields...)
		return true
	})
}

func WarnCtxF(ctx context.Context, format string, args ...interface{}) {
	Logger.Range(func(k, v interface{}) bool {
		v.(logx.GLog).WarnF(format, args...)
		return true
	})
}

func ErrorCtx(ctx context.Context, msg string, fields ...uzap.Field) {
	fields = GetLabels(ctx, fields...)
	Logger.Range(func(k, v interface{}) bool {
		v.(logx.GLog).Error(msg, fields...)
		return true
	})
}

func ErrorCtxF(ctx context.Context, format string, args ...interface{}) {
	Logger.Range(func(k, v interface{}) bool {
		v.(logx.GLog).ErrorF(format, args...)
		return true
	})
}
