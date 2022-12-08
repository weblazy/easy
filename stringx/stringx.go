package stringx

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

var (
	TimeLayout = "2006-01-02 15:04:05"
	ByteSeed   = []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")
)

func ToString(param interface{}) (string, error) {
	resp := ""
	switch param.(type) {
	case int64:
		resp = strconv.FormatInt(param.(int64), 10)
	case int32:
		resp = strconv.FormatInt(param.(int64), 10)
	case int:
		resp = strconv.Itoa(param.(int))
	case float64:
		resp = strconv.FormatFloat(param.(float64), 'f', -1, 64)
	case float32:
		resp = strconv.FormatFloat(param.(float64), 'f', -1, 64)
	case string:
		resp = param.(string)
	case []byte:
		resp = string(param.([]byte))
	case time.Time:
		resp = param.(time.Time).Format(TimeLayout)
	case *time.Time:
		resp = param.(*time.Time).Format(TimeLayout)
	default:
		return resp, fmt.Errorf("%v is not base type", param)
	}
	return resp, nil
}

func SplitN(s string, n int) []string {
	len := len(s)
	var resp []string
	var index, next int

	for len > index {
		next += n
		if len >= next {
			resp = append(resp, s[index:next])
		} else {
			resp = append(resp, s[index:len])
		}
		index = next
	}
	return resp
}

func ToStr(param interface{}) string {
	resp := ""
	switch param.(type) {
	case int64:
		resp = strconv.FormatInt(param.(int64), 10)
	case int32:
		resp = strconv.FormatInt(param.(int64), 10)
	case int:
		resp = strconv.Itoa(param.(int))
	case float64:
		resp = strconv.FormatFloat(param.(float64), 'f', -1, 64)
	case float32:
		resp = strconv.FormatFloat(param.(float64), 'f', -1, 64)
	case string:
		resp = param.(string)
	case []byte:
		resp = string(param.([]byte))
	case time.Time:
		resp = param.(time.Time).Format(TimeLayout)
	case *time.Time:
		resp = param.(*time.Time).Format(TimeLayout)
	default:
		return resp
	}
	return resp
}

func RandomString(len int) string {
	rand.Seed(time.Now().UnixNano())
	resp := make([]byte, len)
	for i := 0; i < len; i++ {
		resp[i] = ByteSeed[rand.Intn(62)]
	}
	return string(resp)
}
