package ctask

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestDo(t *testing.T) {
	type T = int // task type
	type R = int // result type

	type args struct {
		ctx        context.Context
		ctxTimeout time.Duration
		tasks      []T
		executor   func(ctx context.Context, t T) (R, error)
		opts       []DoOpt
	}
	tests := []struct {
		name       string
		args       args
		want       []R
		requireErr require.ErrorAssertionFunc
	}{
		{
			name: "happy path",
			args: args{
				ctx:      context.Background(),
				tasks:    []T{0, 1, 2, 3, 4, 5, 6},
				executor: fibonacci,
				opts:     nil,
			},
			want:       []R{1, 1, 2, 3, 5, 8, 13},
			requireErr: require.NoError,
		},
		{
			name: "empty slice",
			args: args{
				ctx:      context.Background(),
				tasks:    []T{0, 1, 2, 3, 4, 5, 6},
				executor: fibonacci,
				opts:     nil,
			},
			want:       []R{1, 1, 2, 3, 5, 8, 13},
			requireErr: require.NoError,
		},
		{
			name: "error path & ensure tasks after error aren't executed (1000th fibonacci is too slow to be computed)",
			args: args{
				ctx:      context.Background(),
				tasks:    []T{0, 1, 2, 1, -1, 1000},
				executor: fibonacci,
				opts:     []DoOpt{WithWorkerNum(1)},
			},
			want: nil,
			requireErr: func(t require.TestingT, err error, i ...interface{}) {
				require.Equal(t, errors.New("negative"), err)
			},
		},
		{
			name: "slow function with sleeps should run concurrently without context deadline error",
			args: args{
				ctx:        context.Background(),
				ctxTimeout: 50 * time.Millisecond,
				tasks: []T{
					10, 10, 10, 10, 10, 10, 10, 10, 10, 10,
					10, 10, 10, 10, 10, 10, 10, 10, 10, 10,
					10, 10, 10, 10, 10, 10, 10, 10, 10, 10,
					10, 10, 10, 10, 10, 10, 10, 10, 10, 10,
					10, 10, 10, 10, 10, 10, 10, 10, 10, 10,
				},
				executor: func(ctx context.Context, t T) (R, error) {
					time.Sleep(time.Duration(t) * time.Millisecond)
					return 1, nil
				},
				opts: []DoOpt{WithWorkerNum(20)},
			},
			want: []R{
				1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
				1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
				1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
				1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
				1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
			},
			requireErr: require.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := tt.args.ctx
			if tt.args.ctxTimeout > 0 {
				var cancel context.CancelFunc
				ctx, cancel = context.WithTimeout(ctx, tt.args.ctxTimeout)
				defer cancel()
			}
			got, err := Do(ctx, tt.args.tasks, tt.args.executor, tt.args.opts...)
			tt.requireErr(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}

func fibonacci(ctx context.Context, n int) (int, error) {
	if n < 0 {
		return 0, errors.New("negative")
	}
	if n < 2 {
		return 1, nil
	}
	r1, err := fibonacci(ctx, n-1)
	if err != nil {
		return 0, err
	}

	r2, err := fibonacci(ctx, n-2)
	if err != nil {
		return 0, err
	}

	return r1 + r2, nil
}
