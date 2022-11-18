package code_err

import (
	"context"
	"fmt"

	"github.com/weblazy/easy/utils/elog"
)

var (
	ParamsErr  = NewCodeErr(110003, "参数错误")
	TokenErr   = NewCodeErr(110004, "无效Token")
	EncryptErr = NewCodeErr(110022, "加密失败")
	DecryptErr = NewCodeErr(110023, "解密失败")
	SignErr    = NewCodeErr(110024, "签名失败")
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

// 打印log
func ErrLog(ctx context.Context, codeErr *CodeErr, err error) error {
	if _, ok := err.(*CodeErr); ok {
		return err
	}
	elog.ErrorCtx(ctx, "CodeErr", elog.FieldError(err))
	return codeErr
}

// 打印log
func ErrLogf(ctx context.Context, codeErr *CodeErr, format string, a ...interface{}) error {
	elog.ErrorCtx(ctx, fmt.Sprintf(format, a...))
	return codeErr
}
