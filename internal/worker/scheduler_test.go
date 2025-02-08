package worker

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestScheduler(t *testing.T) {
	t.Run("should call process function at least once", func(t *testing.T) {
		// given
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		var callCount int
		process := func(ctx context.Context) error {
			callCount++
			return nil
		}

		// when
		go Scheduler(ctx, 50*time.Millisecond, process)

		time.Sleep(150 * time.Millisecond)
		cancel()

		// then
		assert.GreaterOrEqual(t, callCount, 1)
	})

	t.Run("should stop calling process if context is canceled", func(t *testing.T) {
		// given
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		var callCount int
		process := func(ctx context.Context) error {
			callCount++
			return nil
		}

		// when
		go Scheduler(ctx, 50*time.Millisecond, process)
		time.Sleep(150 * time.Millisecond)

		// then
		assert.Equal(t, 0, callCount, "process should not be called if context is canceled")
	})

	t.Run("should continue processing if function returning error", func(t *testing.T) {
		// given
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		var called bool
		processInt := func(ctx context.Context) error {
			called = true
			return errors.New("error")
		}

		// when
		go Scheduler(ctx, 50*time.Millisecond, processInt)
		time.Sleep(100 * time.Millisecond)
		cancel()

		// then
		assert.True(t, called, "process function should have been called")
	})
}
