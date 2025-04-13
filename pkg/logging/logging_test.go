package logging

import (
	"testing"
)

func TestLogInfo(t *testing.T) {
	// Test logging info message
	LogInfo("Test info message")
}

func TestLogError(t *testing.T) {
	// Test logging error message
	LogError("Test error message")
}
