package eviper

import (
	"os"
	"testing"
)

func TestGetEnvConfig(t *testing.T) {
	v := NewConfig("", "")
	err := os.Setenv("TEST_DEMO", "6666")
	if err != nil {
		t.Error(err)
		return
	}
	s := v.GetEnvConfig("test.demo").String()
	if s != "6666" {
		t.Failed()
	}

	i64 := v.GetEnvConfig("test.demo").Int64()
	if i64 != 6666 {
		t.Failed()
	}
}
