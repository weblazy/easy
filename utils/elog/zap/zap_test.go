package zap

import (
	"testing"

	"github.com/sunmi-OS/gocore/v2/utils/file"
)

func TestName(t *testing.T) {
	deleteLog(file.GetPath()+"/RunTime", 7)
}
