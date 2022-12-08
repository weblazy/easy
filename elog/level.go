package elog

import "fmt"

type CtxLogLevelKey struct{}

// LogLevel log level
type LogLevel int

const (
	// Silent silent log level
	Silent LogLevel = iota + 1
	// Error error log level
	Error
	// Warn warn log level
	Warn
	// Info info log level
	Info
	// Debug debug log level
	Debug
)

// String returns a lower-case ASCII representation of the log level.
func (l LogLevel) String() string {
	switch l {
	case Silent:
		return "silent"
	case Error:
		return "error"
	case Warn:
		return "warn"
	case Info:
		return "info"
	case Debug:
		return "debug"
	default:
		return fmt.Sprintf("Level(%d)", l)
	}
}
