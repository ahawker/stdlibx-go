package stdlib

import (
	"context"
	"sync"
	"time"
)

var (
	// DefaultTaskTimeout is the default timeout for a task to run.
	DefaultTaskTimeout = 30 * time.Second
	// DefaultCancelFn is the default cancel function for a task.
	DefaultCancelFn = func(context.Context) error { return nil }
	// DefaultCancelTimeout is the default timeout for a task cancel fn to run.
	DefaultCancelTimeout = 30 * time.Second
)

// ErrTaskTimeout is returned when a task reaches its timeout cancelled.
var ErrTaskTimeout = Error{
	Code:      "task_timeout",
	Message:   "task reached its timeout and was cancelled",
	Namespace: ErrorNamespaceDefault,
}

// TaskFn is a function that represents a task to be executed.
type TaskFn func(context.Context) error

// TaskCancelFn is a function that represents a task cancellation function.
type TaskCancelFn func(ctx context.Context) error

// TaskConfig for a task execution.
type TaskConfig struct {
	// Timeout is the max duration for the task to run before cancellation.
	Timeout time.Duration
	// Cancel is the function to call when the task is cancelled.
	Cancel TaskCancelFn
	// CancelTimeout is the max duration for the 'TaskCancelFn' to run.
	CancelTimeout time.Duration
}

// WithTaskTimeout sets the timeout for the task.
func WithTaskTimeout(timeout time.Duration) Option[*TaskConfig] {
	return func(o *TaskConfig) error {
		o.Timeout = timeout
		return nil
	}
}

// WithTaskCancel sets the cancel function for the task.
func WithTaskCancel(fn TaskCancelFn) Option[*TaskConfig] {
	return func(o *TaskConfig) error {
		o.Cancel = fn
		return nil
	}
}

// Task executes the given task with the provided context and options.
func Task(ctx context.Context, task TaskFn, options ...Option[*TaskConfig]) error {
	cfg, err := OptionApply(
		&TaskConfig{
			Timeout:       DefaultTaskTimeout,
			Cancel:        DefaultCancelFn,
			CancelTimeout: DefaultCancelTimeout,
		},
		options...,
	)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeoutCause(ctx, cfg.Timeout, ErrTaskTimeout)
	defer cancel()

	var wg sync.WaitGroup
	wg.Add(2)

	eg := NewErrorGroup()
	done := make(chan struct{})

	// Worker.
	go func() {
		defer wg.Done()
		eg.Append(task(ctx))
		close(done)
	}()

	// Watcher.
	go func() {
		defer wg.Done()
		for {
			select {
			case <-done:
				return
			case <-ctx.Done():
				eg.Append(context.Cause(ctx), cfg.Cancel(ctx))
				return
			}
		}
	}()

	wg.Wait()
	return eg.ErrorOrNil()
}
