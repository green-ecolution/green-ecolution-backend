package worker

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type TestWork struct {
	callCount int
}

func (t *TestWork) Do(ctx context.Context) error {
	t.callCount++
	return nil
}

type ErrorWork struct {
	callCount int
}

func (e *ErrorWork) Do(ctx context.Context) error {
	e.callCount++
	return errors.New("intentional error")
}

func TestRunScheduler(t *testing.T) {
	t.Run("should call Do() function at least once", func(t *testing.T) {
		// given
		ctx := context.Background()

		var callCount int
		process := func(ctx context.Context) error {
			callCount++
			return nil
		}

		// when
		scheduler := RunScheduler(ctx, 50*time.Millisecond, process)

		time.Sleep(150 * time.Millisecond)
		scheduler.Stop()

		// then
		assert.GreaterOrEqual(t, callCount, 1)
	})

	t.Run("should stop calling Do() if context is canceled", func(t *testing.T) {
		// given
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		var callCount int
		process := func(ctx context.Context) error {
			if ctx.Err() != nil {
				return ctx.Err()
			}
			callCount++
			return nil
		}

		// when
		RunScheduler(ctx, 50*time.Millisecond, process)
		time.Sleep(150 * time.Millisecond)

		// then
		assert.Equal(t, 0, callCount, "Do() should not be called if context is canceled")
	})

	t.Run("should continue processing if function returning error", func(t *testing.T) {
		// given
		ctx := context.Background()

		var called bool
		processInt := func(ctx context.Context) error {
			called = true
			return errors.New("error")
		}

		// when
		scheduler := RunScheduler(ctx, 50*time.Millisecond, processInt)
		time.Sleep(100 * time.Millisecond)
		scheduler.Stop()

		// then
		assert.True(t, called, "Do() function should have been called")
	})
}

func TestSchedulerWithStruct(t *testing.T) {
	t.Run("should create a scheduler instance and execute Do() at least once", func(t *testing.T) {
		// given
		ctx := context.Background()
		task := &TestWork{}

		// when
		scheduler := NewScheduler(ctx, 50*time.Millisecond, task)
		go scheduler.Run()

		time.Sleep(150 * time.Millisecond)
		scheduler.Stop()

		// then
		assert.GreaterOrEqual(t, task.callCount, 1, "Scheduler should execute Do() at least once")
	})

	t.Run("should stop execution when Stop() is called", func(t *testing.T) {
		// given
		ctx := context.Background()
		task := &TestWork{}

		// when
		scheduler := NewScheduler(ctx, 50*time.Millisecond, task)
		go scheduler.Run()

		time.Sleep(150 * time.Millisecond)
		scheduler.Stop()
		countAfterStop := task.callCount

		time.Sleep(150 * time.Millisecond)

		// then
		assert.Equal(t, countAfterStop, task.callCount, "Scheduler should stop executing Do() after Stop() is called")
	})

	t.Run("should continue execution even if Do() returns an error", func(t *testing.T) {
		// given
		ctx := context.Background()
		task := &ErrorWork{}

		// when
		scheduler := NewScheduler(ctx, 50*time.Millisecond, task)
		go scheduler.Run()

		time.Sleep(150 * time.Millisecond)
		scheduler.Stop()

		// then
		assert.GreaterOrEqual(t, task.callCount, 1, "Scheduler should continue running even if Do() returns an error")
	})

	t.Run("should not execute any Do() calls if started and stopped immediately", func(t *testing.T) {
		// given
		ctx := context.Background()
		task := &TestWork{}

		// when
		scheduler := NewScheduler(ctx, 50*time.Millisecond, task)
		go scheduler.Run()
		scheduler.Stop()

		// then
		assert.Equal(t, 0, task.callCount, "If stopped immediately, Do() should not be called")
	})
}
