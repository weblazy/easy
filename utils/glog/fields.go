package glog

import (
	"fmt"
	"runtime"
	"strconv"
	"time"

	"go.uber.org/zap"
)

const (
	KeyComponent     = "comp"
	KeyError         = "error"
	KeyTraceID       = "trace_id"
	KeySpanID        = "span_id"
	KeyEvent         = "event"
	KeyCost          = "duration"
	KeyMethod        = "method"
	KeyName          = "name"
	KeyAddr          = "addr"
	KeyRestyResponse = "resty_response"
	KeyReq           = "req"
	KeyResp          = "resp"
	KeyDetail        = "detail"
)

// FieldComponent 设置组件.
func FieldComponent(value string) zap.Field {
	return zap.String(KeyComponent, value)
}

// FieldError 设置错误.
func FieldError(err error) zap.Field {
	return zap.Error(err)
}

// FieldTrace 设置 trace id.
func FieldTrace(id string) zap.Field {
	return zap.String(KeyTraceID, id)
}

// FieldSpanID 设置 span id.
func FieldSpanID(id string) zap.Field {
	return zap.String(KeySpanID, id)
}

// FieldEvent 设置 event.
func FieldEvent(event string) zap.Field {
	return zap.String(KeyEvent, event)
}

// FieldCost 设置 duration in ms.
func FieldCost(v time.Duration) zap.Field {
	return zap.Int64(KeyCost, int64(v.Microseconds())/1000)
}

// FieldMethod 设置 method.
func FieldMethod(v string) zap.Field {
	return zap.String(KeyMethod, v)
}

// FieldName 设置 name.
func FieldName(v string) zap.Field {
	return zap.String(KeyName, v)
}

// FieldAddr 设置 addr.
func FieldAddr(v string) zap.Field {
	return zap.String(KeyAddr, v)
}

// FieldReq 设置 req.
func FieldReq(req interface{}) zap.Field {
	return zap.Any(KeyReq, req)
}

// FieldResp 设置 resp.
func FieldResp(resp interface{}) zap.Field {
	return zap.Any(KeyResp, resp)
}

// MakeReqResInfo 以info级别打印行号、配置名、目标地址、耗时、请求数据、响应数据
func MakeReqResInfo(callerSkip int, compName string, addr string, duration time.Duration, req interface{}, reply interface{}) zap.Field {
	_, file, line, _ := runtime.Caller(callerSkip)
	return zap.String(KeyDetail, fmt.Sprintf("%s %s %s %s %s => %s \n", file+":"+strconv.Itoa(line), compName, addr, fmt.Sprintf("[%vms]", float64(duration.Microseconds())/1000), fmt.Sprintf("%v", req), fmt.Sprintf("%v", reply)))
}

// MakeReqResError 以error级别打印行号、配置名、目标地址、耗时、请求数据、响应数据
func MakeReqResError(callerSkip int, compName string, addr string, duration time.Duration, req string, err string) zap.Field {
	_, file, line, _ := runtime.Caller(callerSkip)
	return zap.Error(fmt.Errorf("%s %s %s %s %s => %s ", file+":"+strconv.Itoa(line), compName, addr, fmt.Sprintf("[%vms]", float64(duration.Microseconds())/1000), req, err))
}
