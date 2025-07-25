package logger

import (
	"context"
	"testing"
)

func TestInit(t *testing.T) {
	tests := []struct {
		name     string
		logLevel string
		hasError bool
	}{
		{"Valid info level", "info", false},
		{"Valid debug level", "debug", false},
		{"Valid warn level", "warn", false},
		{"Valid error level", "error", false},
		{"Valid fatal level", "fatal", false},
		{"Invalid log level", "asdasdd--1", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			InitLogger(true)
		})
	}
}

func TestContextualLogging(t *testing.T) {
	// Initialize the logger
	InitLogger(true)

	ctx := context.WithValue(context.Background(), TraceId, "test-request-id")

	tests := []struct {
		name    string
		logFunc func(context.Context, string, ...interface{})
		message string
	}{
		{"Debug", Debug, "Debug message"},
		{"Info", Info, "Info message"},
		{"Error", Error, "Error message"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.logFunc(ctx, tt.message)
		})
	}
}

func TestSync(t *testing.T) {
	InitLogger(true)
	_ = Sync()
}
