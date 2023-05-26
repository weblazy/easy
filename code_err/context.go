package code_err

import (
	"context"

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
	return LogErr(c.Ctx, codeErr, msg, err)
}

// 打印log
func (c *Log) LogField(codeErr *CodeErr, msg string, fields ...zap.Field) *CodeErr {
	return LogField(c.Ctx, codeErr, msg, fields...)
}
