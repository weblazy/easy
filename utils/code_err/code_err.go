package code_err

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
