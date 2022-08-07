package viper

import (
	"os"
	"testing"
)

func TestGetEnvConfig(t *testing.T) {
	err := os.Setenv("TEST_DEMO", "6666")
	if err != nil {
		t.Error(err)
		return
	}
	s := GetEnvConfig("test.demo").String()
	if s != "6666" {
		t.Failed()
	}

	i64 := GetEnvConfig("test.demo").Int64()
	if i64 != 6666 {
		t.Failed()
	}
}
