package ezap

import (
	"testing"

	"github.com/weblazy/easy/filex"
)

func TestName(t *testing.T) {
	deleteLog(filex.GetPath()+"/RunTime", 7)
}
