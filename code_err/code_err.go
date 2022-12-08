package code_err

import (
	"context"

	"github.com/weblazy/easy/utils/elog"
	"go.uber.org/zap"
)

var (
	SystemErr  = NewCodeErr(-1, "系统错误")
	ParamsErr  = NewCodeErr(100001, "参数错误")
	TokenErr   = NewCodeErr(100002, "无效Token")
	EncryptErr = NewCodeErr(100003, "加密失败")
	DecryptErr = NewCodeErr(100004, "解密失败")
	SignErr    = NewCodeErr(100005, "签名失败")
)

type CodeErr struct {
	Code int64
	Msg  string
}

func (err *CodeErr) Error() string {
	return err.Msg
}

func NewCodeErr(code int64, msg string) *CodeErr {
	return &CodeErr{
		Code: code,
		Msg:  msg,
	}
}

func GetCodeErr(err error) *CodeErr {
	if err == nil {
		return nil
	}
	if v, ok := err.(*CodeErr); ok {
		return v
	}
	return SystemErr
}

// 打印msg和err
func LogErr(ctx context.Context, codeErr *CodeErr, msg string, err error) *CodeErr {
	if v, ok := err.(*CodeErr); ok {
		return v
	}
	elog.ErrorCtx(ctx, msg, elog.FieldError(err))
	return codeErr
}

// 打印field
func LogField(ctx context.Context, codeErr *CodeErr, msg string, fields ...zap.Field) *CodeErr {
	elog.ErrorCtx(ctx, msg, fields...)
	return codeErr
}
