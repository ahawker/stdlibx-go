package stdlib

import (
	"context"
	"fmt"
	"github.com/ahawker/stdlibx-go/stdtest"
	"testing"
	"time"
)

func TestTaskRun(t *testing.T) {
	type Got struct {
		task    TaskFn
		options []Option[*TaskConfig]
	}
	stdtest.Table[Got, any]{
		"pass: no timeout with defaults": {
			Got: Got{
				task: func(ctx context.Context) error { return nil },
			},
		},
		"fail: task takes longer than timeout": {
			Got: Got{
				task: func(ctx context.Context) error {
					time.Sleep(500 * time.Millisecond)
					return nil
				},
				options: []Option[*TaskConfig]{
					WithTaskTimeout(100 * time.Millisecond),
				},
			},
			WantErr: ErrTaskTimeout,
		},
		"fail: task takes longer than timeout with cancel that succeeds": {
			Got: Got{
				task: func(ctx context.Context) error {
					time.Sleep(500 * time.Millisecond)
					return nil
				},
				options: []Option[*TaskConfig]{
					WithTaskTimeout(100 * time.Millisecond),
					WithTaskCancel(func(ctx context.Context) error {
						return nil
					}),
				},
			},
			WantErr: ErrTaskTimeout,
		},
		"fail: task takes longer than timeout with cancel that also fails": {
			Got: Got{
				task: func(ctx context.Context) error {
					time.Sleep(500 * time.Millisecond)
					return nil
				},
				options: []Option[*TaskConfig]{
					WithTaskTimeout(100 * time.Millisecond),
					WithTaskCancel(func(ctx context.Context) error {
						return fmt.Errorf("cancel failed")
					}),
				},
			},
			WantErr: ErrTaskTimeout,
		},
	}.Run(t, func(t *stdtest.Test, tc stdtest.Testcase[Got, any]) {
		err := Task(t.Config.Context, tc.Got.task, tc.Got.options...)
		if tc.WantErr != nil {
			t.NotOK(err)
			t.EqualError(err, tc.WantErr)
		} else {
			t.OK(err)
		}
	})
}
