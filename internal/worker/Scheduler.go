package worker

import (
	"context"
	"log/slog"
	"time"
)

func Scheduler(ctx context.Context, interval time.Duration, process func(ctx context.Context) error) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	if ctx.Err() == nil {
		err := process(ctx)
		if err != nil {
			slog.Error("Error during initial process execution", "error", err)
		}
	} else {
		slog.Info("Stopping scheduler before first execution due to canceled context")
		return
	}

	for {
		select {
		case <-ticker.C:
			err := process(ctx)
			if err != nil {
				slog.Error("Error during process execution", "error", err)
			}
		case <-ctx.Done():
			slog.Info("Stopping scheduler")
			return
		}
	}
}
