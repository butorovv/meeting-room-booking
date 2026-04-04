package logger

import (
	"context"
	"log/slog"
	"os"
)

type Logger = *slog.Logger

type contextKey string

const (
	loggerKey    contextKey = "logger"
	requestIDKey contextKey = "request_id"
)

var std Logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))

func ContextWithLogger(ctx context.Context, l Logger) context.Context {
	return context.WithValue(ctx, loggerKey, l)
}

func FromContext(ctx context.Context) Logger {
	if ctx == nil {
		return std
	}
	if v := ctx.Value(loggerKey); v != nil {
		if l, ok := v.(Logger); ok {
			return l
		}
	}
	return std
}

func ContextWithRequestID(ctx context.Context, reqID string) context.Context {
	return context.WithValue(ctx, requestIDKey, reqID)
}

func RequestIDFromContext(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	if v := ctx.Value(requestIDKey); v != nil {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

func Global() Logger {
	return std
}
