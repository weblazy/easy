package env

import (
	"fmt"
	"os"
)

const (
	ReleaseEnv = "onl"
)

var releaseFlag = false //为true时表示线上环境

// GetRunTime 获取当前系统环境
func GetRunTime() string {
	RunTime := os.Getenv("RUN_TIME")
	if RunTime == "" {
		fmt.Println("No RUN_TIME Can't start")
	}
	return RunTime
}

// OnRelease 开启线上环境
func OnRelease() {
	releaseFlag = true
}

// IsRelease 如果是线上环境返回true
func IsRelease() bool {
	return releaseFlag || GetRunTime() == ReleaseEnv
}
