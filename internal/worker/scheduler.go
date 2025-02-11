package worker

import (
	"context"
	"time"

	"github.com/green-ecolution/green-ecolution-backend/internal/logger"
)

type Scheduler struct {
	ctx      context.Context
	cancel   context.CancelFunc
	interval time.Duration
	work     SchedulerWork
}

func NewScheduler(ctx context.Context, interval time.Duration, w SchedulerWork) *Scheduler {
	ctx, cancel := context.WithCancel(ctx)
	return &Scheduler{
		ctx:      ctx,
		cancel:   cancel,
		interval: interval,
		work:     w,
	}
}

type SchedulerWork interface {
	Do(context.Context) error
}

type SchedulerFunc func(context.Context) error

func (s SchedulerFunc) Do(ctx context.Context) error {
	return s(ctx)
}

func RunScheduler(ctx context.Context, interval time.Duration, fn SchedulerFunc) *Scheduler {
	s := NewScheduler(ctx, interval, fn)
	go s.Run()
	return s
}

func (s *Scheduler) Run() {
	log := logger.GetLogger(s.ctx)
	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()

	if err := s.work.Do(s.ctx); err != nil {
		log.Error("error during initial process execution", "error", err)
	}

	for {
		select {
		case <-ticker.C:
			if err := s.work.Do(s.ctx); err != nil {
				log.Error("error during process execution", "error", err)
			}
		case <-s.ctx.Done():
			log.Debug("stopping scheduler")
			return
		}
	}
}

func (s *Scheduler) Stop() {
	s.cancel()
}
