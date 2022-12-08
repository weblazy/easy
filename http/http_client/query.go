package http_client

import (
	"fmt"
	"net/url"

	"github.com/weblazy/easy/utils/stringx"
)

// 编译http参数
func MapToQuery(params map[string]interface{}, urlEncode ...bool) (string, error) {
	if params == nil {
		return "", fmt.Errorf("param is nil")
	}
	v := make(url.Values)
	for key := range params {
		value, err := stringx.ToString(params[key])
		if err != nil {
			return "", nil
		}
		v.Add(key, value)
	}
	encodeStr := v.Encode()
	if len(urlEncode) > 0 && urlEncode[0] {
		return encodeStr, nil
	}
	decodeStr, _ := url.QueryUnescape(encodeStr)
	return decodeStr, nil

}
