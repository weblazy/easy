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

func Infof(format string, a ...interface{}) {
	_, file, line, ok := runtime.Caller(1)
	if ok {
		data, _ := json.Marshal(&Param{
			Time: time.Now().Format("2006-01-02 15:04:05"),
			File: fmt.Sprintf("%s:%d", file, line),
			Data: fmt.Sprintf(format, a...),
		})
		fmt.Printf("%s\n", string(data))
	}
}

func Stack(args ...interface{}) {
	param := &Param{
		Time: time.Now().Format("2006-01-02 15:04:05"),
		File: "",
		Data: args,
	}

	for i := 1; ; i++ {
		pc, file, line, ok := runtime.Caller(i)
		if !ok {
			data, _ := json.Marshal(param)
			fmt.Printf("%s\n", string(data))
			break
		}
		f := runtime.FuncForPC(pc)
		if f.Name() != "runtime.main" && f.Name() != "runtime.goexit" {
			param.File += fmt.Sprintf("%s:%d|------|", file, line)
		}
	}
}
