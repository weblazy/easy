package eviper

import (
	"os"
	"testing"
)

func TestGetEnvConfig(t *testing.T) {
	v := NewViperFromString("")
	err := os.Setenv("TEST_DEMO", "6666")
	if err != nil {
		t.Error(err)
		return
	}
	s := v.GetEnvConfig("test.demo")
	if s != "6666" {
		t.Failed()
	}

}
