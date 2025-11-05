package code_err

import (
	"context"

	"github.com/weblazy/easy/elog"
	"go.uber.org/zap"
)

var (
	SystemErr  = NewCodeErr(-1, "SystemError")
	ParamsErr  = NewCodeErr(100001, "ParamsError")
	TokenErr   = NewCodeErr(100002, "InvalidToken")
	EncryptErr = NewCodeErr(100003, "EncryptionError")
	DecryptErr = NewCodeErr(100004, "DecryptionError")
	SignErr    = NewCodeErr(100005, "SignatureError")
)

type CodeErr struct {
	Code     int64  `json:"code"`
	Msg      string `json:"msg"`
	DebugMsg string `json:"debug_msg"`
}

func (err *CodeErr) Error() string {
	return err.Msg
}

func New(code int64, msg string, debugMsg string) *CodeErr {
	return &CodeErr{
		Code:     code,
		Msg:      msg,
		DebugMsg: debugMsg,
	}
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
func (codeErr *CodeErr) LogErr(ctx context.Context, msg string, err error) *CodeErr {
	if v, ok := err.(*CodeErr); ok {
		return v
	}
	elog.ErrorCtx(elog.AddCtxSkip(ctx, 1), msg, elog.FieldError(err))
	return codeErr.WithDebugMsg(err.Error())
}

// 打印field
func (codeErr *CodeErr) LogField(ctx context.Context, msg string, fields ...zap.Field) *CodeErr {
	elog.ErrorCtx(elog.AddCtxSkip(ctx, 1), msg, fields...)
	return codeErr
}

// debug msg初始化新的CodeErr
func (codeErr *CodeErr) WithDebugMsg(debugMsg string) *CodeErr {
	return New(codeErr.Code, codeErr.Msg, debugMsg)
}

// 设置新的debug msg
func (codeErr *CodeErr) SetDebugMsg(debugMsg string) *CodeErr {
	codeErr.DebugMsg = debugMsg
	return codeErr
}
