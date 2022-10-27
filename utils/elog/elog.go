package elog

import (
	"go.uber.org/zap"
	"golang.org/x/net/context"
)

type CtxFieldKey struct{}
type CtxSkipKey struct{}

const DefaultSkip = 2

func AddCtxSkip(ctx context.Context, skip int) context.Context {
	v, _ := ctx.Value(CtxSkipKey{}).(int)
	return context.WithValue(ctx, CtxSkipKey{}, v+skip)
}

func GetCtxSkip(ctx context.Context) int {
	v, _ := ctx.Value(CtxSkipKey{}).(int)
	return v
}

func DebugCtx(ctx context.Context, msg string, fields ...zap.Field) {
	logLevel, ok := ctx.Value(CtxSkipKey{}).(LogLevel)
	if ok && logLevel < Debug {
		return
	}
	fields = MergeCtxFields(ctx, fields...)
	logger := GetLoggerFromCtx(ctx)
	if logger != nil {
		logger.DebugCtx(ctx, msg, fields...)
	}
}

func InfoCtx(ctx context.Context, msg string, fields ...zap.Field) {
	logLevel, ok := ctx.Value(CtxSkipKey{}).(LogLevel)
	if ok && logLevel < Info {
		return
	}
	fields = MergeCtxFields(ctx, fields...)

	logger := GetLoggerFromCtx(ctx)
	if logger != nil {
		logger.InfoCtx(ctx, msg, fields...)
	}
}

func WarnCtx(ctx context.Context, msg string, fields ...zap.Field) {
	logLevel, ok := ctx.Value(CtxSkipKey{}).(LogLevel)
	if ok && logLevel < Warn {
		return
	}
	fields = MergeCtxFields(ctx, fields...)
	logger := GetLoggerFromCtx(ctx)
	if logger != nil {
		logger.WarnCtx(ctx, msg, fields...)
	}
}

func ErrorCtx(ctx context.Context, msg string, fields ...zap.Field) {
	logLevel, ok := ctx.Value(CtxSkipKey{}).(LogLevel)
	if ok && logLevel < Error {
		return
	}
	fields = MergeCtxFields(ctx, fields...)
	logger := GetLoggerFromCtx(ctx)
	if logger != nil {
		logger.ErrorCtx(ctx, msg, fields...)
	}
}
