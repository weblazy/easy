package ezap

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Config struct {
	Logfile   *os.File
	ZapConfig zap.Config
	MaxAge    int64 //定期清理日志文件,日志保留天数
}

// DefaultConfig default config ...
func DefaultConfig() *Config {
	zapConfig := zap.NewProductionConfig()
	// zapConfig.EncoderConfig.TimeKey = zapcore.OmitKey
	// zapConfig.EncoderConfig.LevelKey = zapcore.OmitKey
	// zapConfig.EncoderConfig.NameKey = zapcore.OmitKey
	// zapConfig.EncoderConfig.CallerKey = zapcore.OmitKey
	// zapConfig.EncoderConfig.FunctionKey = zapcore.OmitKey
	// zapConfig.EncoderConfig.MessageKey = "msg"
	zapConfig.EncoderConfig.StacktraceKey = zapcore.OmitKey
	// zapConfig.EncoderConfig.LineEnding = zapcore.DefaultLineEnding
	// zapConfig.EncoderConfig.EncodeLevel = zapcore.LowercaseLevelEncoder
	zapConfig.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	// zapConfig.EncoderConfig.EncodeDuration = zapcore.SecondsDurationEncoder
	// zapConfig.EncoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	zapConfig.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	return &Config{
		MaxAge:    7,
		ZapConfig: zapConfig,
	}
}
