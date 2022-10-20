package logx

import "go.uber.org/zap"

type GLog interface {
	Info(msg string, fields ...zap.Field)
	InfoF(format string, args ...interface{})
	Debug(msg string, fields ...zap.Field)
	DebugF(format string, args ...interface{})
	Warn(msg string, fields ...zap.Field)
	WarnF(format string, args ...interface{})
	Error(msg string, fields ...zap.Field)
	ErrorF(format string, args ...interface{})
}
