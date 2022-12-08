package run

import (
	"context"
	"runtime/debug"
	"time"

	"emperror.dev/errors"
	"github.com/weblazy/easy/utils/elog"
	"github.com/weblazy/easy/utils/timex"
	"go.uber.org/zap"
)

func DaemonRun(interval time.Duration, f func(), daemon func()) {
	ticker := timex.NewRealTicker(interval)
	stopChannel := make(chan struct{})
	defer func() {
		stopChannel <- struct{}{}
	}()
	go func() {
		for {
			select {
			case <-ticker.Chan():
				daemon()
			case <-stopChannel:
				ticker.Stop()
				return
			}
		}
	}()
	f()
}

// RunSafeWrap wrapper func () error with Recover
func RunSafeWrap(ctx context.Context, fn func() error) (err error) {
	defer func() {
		if p := recover(); p != nil {
			elog.ErrorCtx(ctx, "panic", zap.Any("err", p), zap.String("stack", string(debug.Stack())))
			err = errors.Errorf("panic: %v", p)
		}
	}()

	err = fn()

	return
}
