package code_err

import (
	"context"

	"go.uber.org/zap"
)

type SvcContext struct {
	Ctx context.Context
}

func NewSvcContext(ctx context.Context) *SvcContext {
	return &SvcContext{Ctx: ctx}
}

// 打印log
func (c *SvcContext) LogErr(codeErr *CodeErr, msg string, err error) *CodeErr {
	return LogErr(c.Ctx, codeErr, msg, err)
}

// 打印log
func (c *SvcContext) LogField(codeErr *CodeErr, msg string, fields ...zap.Field) *CodeErr {
	return LogField(c.Ctx, codeErr, msg, fields...)
}
