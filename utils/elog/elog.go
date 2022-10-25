package elog

import (
	uzap "go.uber.org/zap"
	"golang.org/x/net/context"
)

type LogConfCtxKey struct{}

type LogConf struct {
	Name   string
	Labels []uzap.Field
}

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

func DebugCtx(ctx context.Context, msg string, fields ...uzap.Field) {
	logLevel, ok := ctx.Value(LogLevelCtxKey{}).(LogLevel)
	if ok && logLevel < Debug {
		return
	}
	fields = MergeLabels(ctx, fields...)
	logger := GetLoggerFromCtx(ctx)
	if logger != nil {
		logger.DebugCtx(ctx, msg, fields...)
	}
}

func InfoCtx(ctx context.Context, msg string, fields ...uzap.Field) {
	logLevel, ok := ctx.Value(LogLevelCtxKey{}).(LogLevel)
	if ok && logLevel < Info {
		return
	}
	fields = MergeLabels(ctx, fields...)

	logger := GetLoggerFromCtx(ctx)
	if logger != nil {
		logger.InfoCtx(ctx, msg, fields...)
	}
}

func WarnCtx(ctx context.Context, msg string, fields ...uzap.Field) {
	logLevel, ok := ctx.Value(LogLevelCtxKey{}).(LogLevel)
	if ok && logLevel < Warn {
		return
	}
	fields = MergeLabels(ctx, fields...)
	logger := GetLoggerFromCtx(ctx)
	if logger != nil {
		logger.WarnCtx(ctx, msg, fields...)
	}
}

func ErrorCtx(ctx context.Context, msg string, fields ...uzap.Field) {
	logLevel, ok := ctx.Value(LogLevelCtxKey{}).(LogLevel)
	if ok && logLevel < Error {
		return
	}
	fields = MergeLabels(ctx, fields...)
	logger := GetLoggerFromCtx(ctx)
	if logger != nil {
		logger.ErrorCtx(ctx, msg, fields...)
	}
}
