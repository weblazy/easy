package logx

import (
	"encoding/json"
	"fmt"
	"runtime"
	"time"
)

type (
	Param struct {
		Time string      `json:"time"`
		File string      `json:"file"`
		Data interface{} `json:"data"`
	}
)

func Info(args ...interface{}) {
	_, file, line, ok := runtime.Caller(1)
	if ok {
		data, _ := json.Marshal(&Param{
			Time: time.Now().Format("2006-01-02 15:04:05"),
			File: fmt.Sprintf("%s:%d", file, line),
			Data: args,
		})
		fmt.Printf("%s\n", string(data))
	}
}
