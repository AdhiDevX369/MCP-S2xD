package logger

import (
	"context"
	"log/slog"
	"os"
)

type contextKey string

const RequestIDKey contextKey = "request_id"

var defaultLogger *slog.Logger

func Init(level, format string) {
	var handler slog.Handler
	opts := &slog.HandlerOptions{
		Level: parseLevel(level),
	}

	if format == "json" {
		handler = slog.NewJSONHandler(os.Stdout, opts)
	} else {
		handler = slog.NewTextHandler(os.Stdout, opts)
	}

	defaultLogger = slog.New(handler)
	slog.SetDefault(defaultLogger)
}

func parseLevel(level string) slog.Level {
	switch level {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

func WithRequestID(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, RequestIDKey, requestID)
}

func GetRequestID(ctx context.Context) string {
	if id, ok := ctx.Value(RequestIDKey).(string); ok {
		return id
	}
	return ""
}

func Info(ctx context.Context, msg string, args ...any) {
	args = appendRequestID(ctx, args)
	slog.Info(msg, args...)
}

func Error(ctx context.Context, msg string, args ...any) {
	args = appendRequestID(ctx, args)
	slog.Error(msg, args...)
}

func Debug(ctx context.Context, msg string, args ...any) {
	args = appendRequestID(ctx, args)
	slog.Debug(msg, args...)
}

func Warn(ctx context.Context, msg string, args ...any) {
	args = appendRequestID(ctx, args)
	slog.Warn(msg, args...)
}

func appendRequestID(ctx context.Context, args []any) []any {
	if id := GetRequestID(ctx); id != "" {
		return append(args, "request_id", id)
	}
	return args
}
