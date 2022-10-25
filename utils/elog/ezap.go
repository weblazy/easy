package elog

import (
	"log"

	"github.com/weblazy/easy/utils/elog/ezap"
	"go.uber.org/zap"
)

var (
	DefaultLogger *ezap.Ezap
)

const (
	Ezap = "ezap"
)

func init() {
	var err error
	cfg := ezap.DefaultConfig()
	l, err := cfg.ZapConfig.Build(zap.AddCallerSkip(2))
	if err != nil {
		log.Printf("l.initZap(),err:%+v", err)
		return
	}
	DefaultLogger = &ezap.Ezap{
		Logger: l,
		Config: cfg,
	}
}
