package logger

import "context"

type Logger interface {
	Debug(ctx context.Context, format string, args ...interface{})
	Info(ctx context.Context, format string, args ...interface{})
	Warn(ctx context.Context, format string, args ...interface{})
	Error(ctx context.Context, format string, args ...interface{})
}

type contextKey struct {
	name string
}

var loggerKey = &contextKey{"logger"}
