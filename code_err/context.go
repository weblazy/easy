package code_err

import (
	"context"

	"github.com/weblazy/easy/elog"
	"go.uber.org/zap"
)

type Log struct {
	Ctx context.Context
}

func NewLog(ctx context.Context) *Log {
	return &Log{Ctx: ctx}
}

// 打印log
func (c *Log) LogErr(codeErr *CodeErr, msg string, err error) *CodeErr {
	return codeErr.LogErr(elog.AddCtxSkip(c.Ctx, 1), msg, err)
}

// 打印log
func (c *Log) LogField(codeErr *CodeErr, msg string, fields ...zap.Field) *CodeErr {
	return codeErr.LogField(elog.AddCtxSkip(c.Ctx, 1), msg, fields...)
}
