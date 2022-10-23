package elog

import (
	uzap "go.uber.org/zap"
	"golang.org/x/net/context"

	"github.com/weblazy/easy/utils/elog/logx"
)

type LogConfCtxKey struct{}

type LogConf struct {
	Name   string
	Labels []uzap.Field
}

var defaultLogConf = &LogConf{}

func GetContextLog(ctx context.Context) *LogConf {
	if v, ok := ctx.Value(LogConfCtxKey{}).(*LogConf); ok {
		return v
	} else {
		return &LogConf{}
	}
}

func SetContextLog(ctx context.Context, log *LogConf) context.Context {
	newCtx := context.WithValue(ctx, LogConfCtxKey{}, log)
	return newCtx
}

func GetLabels(ctx context.Context) []uzap.Field {
	logConf := GetContextLog(ctx)
	llen := len(logConf.Labels)
	newFields := make([]uzap.Field, llen)
	copy(newFields, logConf.Labels)
	return newFields
}

func SetLabels(ctx context.Context, fields ...uzap.Field) []uzap.Field {
	logConf := GetContextLog(ctx)
	logConf.Labels = append(logConf.Labels, fields...)
	return logConf.Labels
}

func MergeLabels(ctx context.Context, fields ...uzap.Field) []uzap.Field {
	logConf := GetContextLog(ctx)
	llen := len(logConf.Labels)
	flen := len(fields)
	newFields := make([]uzap.Field, llen+flen)
	copy(newFields, logConf.Labels)
	copy(newFields[llen:], fields)
	return newFields
}

func InfoCtx(ctx context.Context, msg string, fields ...uzap.Field) {
	fields = MergeLabels(ctx, fields...)
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
	fields = MergeLabels(ctx, fields...)
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
	fields = MergeLabels(ctx, fields...)
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
	fields = MergeLabels(ctx, fields...)
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