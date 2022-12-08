package elog

import (
	"context"
	"testing"

	"go.uber.org/zap"
)

func TestLog(t *testing.T) {
	s := struct {
		Name string
		Age  int
	}{
		Name: "Jerry",
		Age:  18,
	}
	// zap log
	ctx := context.Background()
	DebugCtx(ctx, "zap debug")
	InfoCtx(ctx, "", zap.Any("obj", s))
	WarnCtx(ctx, "zap warn")
	ErrorCtx(ctx, "zap error")

	ctx = context.WithValue(ctx, CtxSkipKey{}, Info)
	// ctx = AddCtxSkip(ctx, 2)
	ctx = AppendCtxFields(ctx, zap.String("name", "lazy"))
	DebugCtx(ctx, "zap debug")
	InfoCtx(ctx, "", zap.Any("obj", s))
	WarnCtx(ctx, "zap warn")
	ErrorCtx(ctx, "zap error")

	// logger := ezap.NewFileEzap("test1")
	// loggerName := "test"
	// SetLogger(loggerName, logger)
	// ctx = SetLogerName(ctx, loggerName)
	// DebugCtx(ctx, "zap debug")
	// InfoCtx(ctx, "", zap.Any("obj", s))
	// WarnCtx(ctx, "zap warn")
	// ErrorCtx(ctx, "zap error")

}
