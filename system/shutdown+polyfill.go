//go:build windows
// +build windows

package system

import "time"

func AddShutdownListener(fn func()) func() {
	return nil
}

func AddWrapUpListener(fn func()) func() {
	return nil
}

func SetTimeoutToForceQuit(duration time.Duration) {
}
