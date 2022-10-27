package ezap

import (
	"log"
)

func NewConsoleEzap() *Ezap {

	var err error
	cfg := DefaultConfig()
	l, err := cfg.ZapConfig.Build()
	if err != nil {
		log.Printf("l.initZap(),err:%+v", err)
		return nil
	}
	return &Ezap{
		Logger: l,
		Config: cfg,
	}
}
