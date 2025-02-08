package worker

import (
	"context"
	"github.com/green-ecolution/green-ecolution-backend/internal/logger"
	"log/slog"
	"time"
)

func Scheduler(ctx context.Context, interval time.Duration, process func(ctx context.Context) error) {
	log := logger.GetLogger(ctx)
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	if ctx.Err() == nil {
		err := process(ctx)
		if err != nil {
			log.Error("error during initial process execution", "error", err)
		}
	} else {
		log.Debug("stopping scheduler before first execution due to canceled context")
		return
	}

	for {
		select {
		case <-ticker.C:
			err := process(ctx)
			if err != nil {
				slog.Error("error during process execution", "error", err)
			}
		case <-ctx.Done():
			slog.Debug("stopping scheduler")
			return
		}
	}
}
