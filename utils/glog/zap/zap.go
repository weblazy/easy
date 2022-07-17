package zap

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/weblazy/easy/utils/glog/logx"

	"github.com/sunmi-OS/gocore/v2/conf/viper"
	"github.com/sunmi-OS/gocore/v2/utils/file"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	LogLevelInfo  = "info"
	LogLevelDebug = "debug"
	LogLevelWarn  = "warn"
	LogLevelError = "error"
)

var (
	LocalLogger *zap.Logger
	logfile     *os.File
	cfg         zap.Config
)

func init() {
	var err error
	cfg = zap.NewProductionConfig()
	cfg.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	l, err := cfg.Build(zap.AddCallerSkip(2))
	if err != nil {
		log.Printf("l.initZap(),err:%+v", err)
		return
	}
	LocalLogger = l
}

func SetLogLevel(logLevel string) {
	switch logLevel {
	case LogLevelDebug:
		cfg.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	case LogLevelInfo:
		cfg.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	case LogLevelWarn:
		cfg.Level = zap.NewAtomicLevelAt(zap.WarnLevel)
	case LogLevelError:
		cfg.Level = zap.NewAtomicLevelAt(zap.ErrorLevel)
	default:
		cfg.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	}

	logger, err := cfg.Build(zap.AddCallerSkip(1))
	if err != nil {
		log.Printf("l.initZap(),err:%+v.\n", err)
		return
	}
	LocalLogger = logger
}

func InitFileLog(logPath ...string) {
	var (
		err  error
		path = file.GetPath() + "/Runtime"
	)
	if len(logPath) == 1 {
		path = logPath[0]
	}

	if !file.CheckDir(path) {
		if err := file.MkdirDir(path); err != nil {
			log.Printf("l.initZap(),err:%+v.\n", err)
		}
	}

	filename := path + "/" + time.Now().Format("2006-01-02") + ".log"
	logfile, err = os.OpenFile(filename, os.O_RDWR|os.O_APPEND, 0666)
	if err != nil {
		logfile, err = os.Create(filename)
		if err != nil {
			log.Println(err)
		}
	}
	cfg.OutputPaths = []string{filename, "stdout"}
	cfg.ErrorOutputPaths = []string{filename, "stderr"}
	cfg.EncoderConfig = zapcore.EncoderConfig{
		TimeKey:        zapcore.OmitKey,
		LevelKey:       zapcore.OmitKey,
		NameKey:        zapcore.OmitKey,
		CallerKey:      zapcore.OmitKey,
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "msg",
		StacktraceKey:  zapcore.OmitKey,
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.EpochTimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
	SetLogLevel(viper.GetEnvConfig("log.level").String())
	viper.C.SetDefault("log.saveDays", "3")
	go updateLogFile(path)
}

// updateLogFile 检测是否跨天了,把记录记录到新的文件目录中
func updateLogFile(logPath string) {
	var err error
	saveDays := viper.GetEnvConfig("log.saveDays").Float64()
	if logPath == "" {
		logPath = file.GetPath() + "/Runtime"
	}
	for {
		now := time.Now()
		// 计算下一个零点
		next := now.Add(time.Hour * 24)
		next = time.Date(next.Year(), next.Month(), next.Day(), 0, 0, 0, 0, next.Location())
		t := time.NewTimer(next.Sub(now))
		<-t.C
		//以下为定时执行的操作
		logfile.Close()
		go deleteLog(logPath, saveDays)
		filename := logPath + "/" + time.Now().Format("2006-01-02") + ".log"
		logfile, err = os.Create(filename)
		if err != nil {
			log.Println(err)
		}
		cfg.ErrorOutputPaths = []string{filename, "stderr"}
		cfg.OutputPaths = []string{filename, "stdout"}
		l, err := cfg.Build()
		if err != nil {
			log.Println(err)
			continue
		}
		LocalLogger = l
	}
}

// deleteLog 删除修改时间在saveDays天前的文件
func deleteLog(source string, saveDays float64) {
	err := filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() && !strings.HasSuffix(info.Name(), ".log") {
			return nil
		}
		t := time.Since(info.ModTime()).Hours()
		if t >= (saveDays-1)*24 {
			e := os.Remove(path)
			if e != nil {
				log.Println(e)
			}
		}
		return err
	})
	if err != nil {
		return
	}
}

// Zap 将文件输出到终端或者文件
type Zap struct {
	logx.GLog
}

func (*Zap) Info(msg string, fields ...zap.Field) {
	LocalLogger.Info(msg, fields...)
}

func (*Zap) InfoF(format string, args ...interface{}) {
	LocalLogger.Info(fmt.Sprintf(format, args...))
}

func (*Zap) Debug(msg string, fields ...zap.Field) {
	LocalLogger.Debug(msg, fields...)
}

func (*Zap) DebugF(format string, args ...interface{}) {
	LocalLogger.Debug(fmt.Sprintf(format, args...))
}

func (*Zap) Warn(msg string, fields ...zap.Field) {
	LocalLogger.Warn(msg, fields...)
}

func (*Zap) WarnF(format string, args ...interface{}) {
	LocalLogger.Warn(fmt.Sprintf(format, args...))
}

func (*Zap) Error(msg string, fields ...zap.Field) {
	LocalLogger.Error(msg, fields...)
}

func (*Zap) ErrorF(format string, args ...interface{}) {
	LocalLogger.Error(fmt.Sprintf(format, args...))
}
