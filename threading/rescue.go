package threading

import "github.com/weblazy/easy/utils/logx"

func Rescue(cleanups ...func()) {
	for _, cleanup := range cleanups {
		cleanup()
	}

	if p := recover(); p != nil {
		logx.Stack(p)
	}
}
