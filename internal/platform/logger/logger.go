package logger

import (
	"context"
	"log/slog"
	"os"
)

type ctxKey string

const RequestIDKey ctxKey = "request_id"

var Log *slog.Logger

func Init(env string) {
	var handler slog.Handler
	opts := &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}

	if env == "production" {
		handler = slog.NewJSONHandler(os.Stdout, opts)
	} else {
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		})
	}

	Log = slog.New(handler)
	slog.SetDefault(Log)
}

func FromContext(ctx context.Context) *slog.Logger {
	if ctx == nil {
		return Log
	}
	if reqID, ok := ctx.Value(RequestIDKey).(string); ok && reqID != "" {
		return Log.With("request_id", reqID)
	}
	return Log
}
