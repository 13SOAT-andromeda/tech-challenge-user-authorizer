package utils

import (
	"testing"
)

func TestLoggers(t *testing.T) {
	if InfoLogger == nil {
		t.Error("InfoLogger should be initialized")
	}
	if ErrorLogger == nil {
		t.Error("ErrorLogger should be initialized")
	}
}
