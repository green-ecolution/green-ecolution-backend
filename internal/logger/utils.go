package logger

import (
	"context"
	"log/slog"
)

func GetLogger(ctx context.Context) *slog.Logger {
	return ctx.Value("logger").(*slog.Logger)
}
