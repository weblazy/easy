package service

type Response struct {
	Data     interface{} `json:"data"`
	Code     int64       `json:"code"`
	Msg      string      `json:"msg"`
	DebugMsg string      `json:"debug_msg"`
}

var defaultResponse Response

func init() {
	defaultResponse = Response{
		Data:     nil,
		Code:     1,
		Msg:      "",
		DebugMsg: "",
	}
}

// NewResponse 获取默认返回内容
func NewResponse() Response {
	return defaultResponse
}

// SetDefaultCode 设置默认返回code码
func SetDefaultCode(code int64) {
	defaultResponse.Code = code
}

// SetDefaultData 设置默认返回data内容
func SetDefaultData(data interface{}) {
	defaultResponse.Data = data
}

// SetDefaultMsg 设置默认返回msg内容
func SetDefaultMsg(msg string) {
	defaultResponse.Msg = msg
}
