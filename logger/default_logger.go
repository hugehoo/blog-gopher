package logger

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"
)

const (
	DEBUG = "DEBUG"
	INFO  = "INFO"
	WARN  = "WARN"
	ERROR = "ERROR"
)

type DefaultLogger struct {
	serviceName string
}

type LogEntry struct {
	Timestamp   string      `json:"timestamp"`
	Level       string      `json:"level"`
	RequestID   interface{} `json:"request_id"`
	UserID      interface{} `json:"user_id"`
	ServiceName string      `json:"service_name"`
	Message     string      `json:"message"`
}

func NewLogger(serviceName string) *DefaultLogger {
	return &DefaultLogger{
		serviceName: serviceName,
	}
}

func (l *DefaultLogger) Debug(ctx context.Context, format string, args ...interface{}) {
	l.log(ctx, DEBUG, format, args...)
}

func (l *DefaultLogger) Info(ctx context.Context, format string, args ...interface{}) {
	l.log(ctx, INFO, format, args...)
}

func (l *DefaultLogger) Warn(ctx context.Context, format string, args ...interface{}) {
	l.log(ctx, WARN, format, args...)
}

func (l *DefaultLogger) Error(ctx context.Context, format string, args ...interface{}) {
	l.log(ctx, ERROR, format, args...)
}

func (l *DefaultLogger) log(ctx context.Context, level string, format string, args ...interface{}) {
	entry := LogEntry{
		Timestamp:   time.Now().Format(time.RFC3339),
		Level:       level,
		RequestID:   ctx.Value("request_id"),
		UserID:      ctx.Value("user_id"),
		ServiceName: l.serviceName,
		Message:     fmt.Sprintf(format, args...),
	}

	jsonBytes, err := json.Marshal(entry)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to marshal log entry: %v\n", err)
		return
	}
	fmt.Println(string(jsonBytes))
}
