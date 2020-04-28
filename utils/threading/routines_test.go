package threading

import (
	"io/ioutil"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRoutineId(t *testing.T) {
	assert.True(t, RoutineId() > 0)
}

func TestRunSafe(t *testing.T) {
	log.SetOutput(ioutil.Discard)

	i := 0

	defer func() {
		assert.Equal(t, 1, i)
	}()

	ch := make(chan struct{})
	go RunSafe(func() {
		defer func() {
			ch <- struct{}{}
		}()

		panic("panic")
	})

	<-ch
	i++
}
