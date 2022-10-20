package elog

import (
	"errors"
	"fmt"
	"testing"

	"github.com/weblazy/easy/utils/elog/zap"
)

func TestLog(t *testing.T) {
	s := struct {
		Name string
		Age  int
	}{
		Name: "Jerry",
		Age:  18,
	}
	// zap log
	InfoF("%+v", s)
	Debug("zap debug")
	Warn("zap warn")
	Error("zap error")
	ErrorF("s.dao.PartnerById(%d),err:%+v", 10086, errors.New("不存在此id"))
	ErrorF("s.dao.CreateOrder(%+v),err:%+v", s, errors.New("创建订单失败"))

	fmt.Println("")

	zap.SetLogLevel(zap.LogLevelWarn)
	InfoF("%+v", s)
	Debug("zap debug")
	Warn("zap warn")
	Error("zap error")

	fmt.Println("")

	zap.InitFileLog()
	Debug("zap debug")
	Warn("zap warn")
	Error("zap error")

}
