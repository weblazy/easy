package glog

import (
	"sync"

	"github.com/weblazy/easy/utils/glog/zap"

	uzap "go.uber.org/zap"

	"github.com/weblazy/easy/utils/glog/logx"
)

var (
	Logger sync.Map
)

//  默认加入zap组件
func init() {
	Logger.Store("zap", &zap.Zap{})
}

// SetLogger 设置日志打印实例,选择输出到文件,终端,阿里云日志等
func SetLogger(name string, logger logx.GLog) {
	Logger.Store(name, logger)
}

// DelLogger 删除日志插件
func DelLogger(name string) {
	Logger.Delete(name)
}

func Info(msg string, fields ...uzap.Field) {
	Logger.Range(func(k, v interface{}) bool {
		v.(logx.GLog).Info(msg, fields...)
		return true
	})
}

func InfoF(format string, args ...interface{}) {
	Logger.Range(func(k, v interface{}) bool {
		v.(logx.GLog).InfoF(format, args...)
		return true
	})
}

func Debug(msg string, fields ...uzap.Field) {
	Logger.Range(func(k, v interface{}) bool {
		v.(logx.GLog).Debug(msg, fields...)
		return true
	})
}

func DebugF(format string, args ...interface{}) {
	Logger.Range(func(k, v interface{}) bool {
		v.(logx.GLog).DebugF(format, args...)
		return true
	})
}

func Warn(msg string, fields ...uzap.Field) {
	Logger.Range(func(k, v interface{}) bool {
		v.(logx.GLog).Warn(msg, fields...)
		return true
	})
}

func WarnF(format string, args ...interface{}) {
	Logger.Range(func(k, v interface{}) bool {
		v.(logx.GLog).WarnF(format, args...)
		return true
	})
}

func Error(msg string, fields ...uzap.Field) {
	Logger.Range(func(k, v interface{}) bool {
		v.(logx.GLog).Error(msg, fields...)
		return true
	})
}

func ErrorF(format string, args ...interface{}) {
	Logger.Range(func(k, v interface{}) bool {
		v.(logx.GLog).ErrorF(format, args...)
		return true
	})
}
