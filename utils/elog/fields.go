package elog

import (
	"fmt"
	"runtime"
	"strconv"
	"time"

	"go.uber.org/zap"
	"golang.org/x/net/context"
)

type CtxFieldsKey struct{}

const (
	KeyComponent     = "comp"
	KeyError         = "error"
	KeyTraceID       = "trace_id"
	KeySpanID        = "span_id"
	KeyEvent         = "event"
	KeyCost          = "duration"
	KeyDuration      = "duration"
	KeyMethod        = "method"
	KeyName          = "name"
	KeyAddr          = "addr"
	KeyRestyResponse = "resty_response"
	KeyReq           = "req"
	KeyResp          = "resp"
	KeyDetail        = "detail"
	KeySlow          = "slow"
)

// FieldComponent 设置组件.
func FieldComponent(value string) zap.Field {
	return zap.String(KeyComponent, value)
}

// FieldError 设置错误.
func FieldError(err error) zap.Field {
	return zap.Error(err)
}

// FieldSlow 慢操作.
func FieldSlow(isSlow bool) zap.Field {
	return zap.Bool(KeySlow, isSlow)
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

// FieldCost 设置 duration in ms.
func FieldDuration(v time.Duration) zap.Field {
	return zap.Int64(KeyDuration, int64(v.Microseconds())/1000)
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

func GetCtxFields(ctx context.Context) []zap.Field {
	fields, _ := ctx.Value(CtxFieldKey{}).([]zap.Field)
	return fields
}

func AppendCtxFields(ctx context.Context, fields ...zap.Field) context.Context {
	if ctxFields, ok := ctx.Value(CtxFieldKey{}).([]zap.Field); ok {
		llen := len(ctxFields)
		flen := len(fields)
		newFields := make([]zap.Field, llen+flen)
		copy(newFields, ctxFields)
		copy(newFields[llen:], fields)
		return context.WithValue(ctx, CtxFieldKey{}, newFields)
	}
	return context.WithValue(ctx, CtxFieldKey{}, fields)
}

func MergeCtxFields(ctx context.Context, fields ...zap.Field) []zap.Field {
	skip := GetCtxSkip(ctx)
	_, file, line, _ := runtime.Caller(DefaultSkip + skip)
	skipField := zap.String("caller", fmt.Sprintf("%s:%d", file, line))
	flen := len(fields)
	// 获取ctx中的field
	if ctxFields, ok := ctx.Value(CtxFieldKey{}).([]zap.Field); ok {
		llen := len(ctxFields)

		newFields := make([]zap.Field, llen+flen)
		copy(newFields, ctxFields)
		copy(newFields[llen:], fields)
		newFields = append(newFields, skipField)
		return newFields
	}
	newFields := make([]zap.Field, flen)
	copy(newFields, fields)
	fields = append(fields, skipField)
	return fields
}
