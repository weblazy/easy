package eerror

import (
	"strconv"

	"github.com/weblazy/easy/set"
)

const (
	UnknownErrCode = "9999"
	SuccessCode    = "20000"
)

var (
	DefaultSuccessCodes = []string{SuccessCode}
)

type commonErrResp interface {
	GetError() error
}

type codeMsgResp interface {
	GetCode() int64
	GetMsg() string
}

type codeMsgStringResp interface {
	GetCode() string
	GetMessage() string
}

type retCodeMsgResp interface {
	GetRetCode() int32
	GetRetMsg() string
}

func ExtractBizCode(successCodes []string) func(resp interface{}, err error) (string, bool) {
	sc := successCodes
	if len(sc) == 0 {
		sc = DefaultSuccessCodes
	}
	scs := set.NewStringSet()
	scs.BatchAdd(sc...)

	replacer := func(s string) string {
		if scs.Has(s) {
			return SuccessCode
		}
		return s
	}

	return func(resp interface{}, err error) (string, bool) {
		// if err != nil {
		// commonErr, ok := FromErrorIsDetail(err)
		// // 1. rich error
		// if ok {
		// 	return commonErr.BizCode, true
		// }
		// non biz error
		// 	return "", false
		// }

		// // 2. 内嵌 commonError
		// if cer, ok := resp.(commonErrResp); ok {
		// 	// 内嵌 error nil 当做成功处理
		// 	if cer.GetError() == nil {
		// 		return SuccessCode, true
		// 	}
		// 	return replacer(cer.GetError().GetBizCode()), true
		// }

		// 3. 内嵌 code msg
		if cmr, ok := resp.(codeMsgResp); ok {
			return replacer(strconv.Itoa(int(cmr.GetCode()))), true
		}

		// 4. 内嵌 code msg string
		if cmr, ok := resp.(codeMsgStringResp); ok {
			return replacer(cmr.GetCode()), true
		}

		// 5. 内嵌 ret_code ret_msg
		if cmr, ok := resp.(retCodeMsgResp); ok {
			return replacer(strconv.Itoa(int(cmr.GetRetCode()))), true
		}

		// response 不符合上面任何标准, 且 err 为 nil
		// 当做正常请求计算
		if err == nil {
			return SuccessCode, true
		}

		// 不属于以上任意一种
		return "", false
	}
}
