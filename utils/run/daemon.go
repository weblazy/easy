package run

import (
	"time"

	"github.com/weblazy/easy/utils/timex"
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
