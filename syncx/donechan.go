package syncx

import (
	"sync"

	"github.com/weblazy/easy/utils/lang"
)

type DoneChan struct {
	done chan lang.PlaceholderType
	once sync.Once
}

func NewDoneChan() *DoneChan {
	return &DoneChan{
		done: make(chan lang.PlaceholderType),
	}
}

func (dc *DoneChan) Close() {
	dc.once.Do(func() {
		close(dc.done)
	})
}

func (dc *DoneChan) Done() chan lang.PlaceholderType {
	return dc.done
}
