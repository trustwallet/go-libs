package ctask

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDo(t *testing.T) {
	type T = int // task type
	type R = int // result type

	type args struct {
		ctx      context.Context
		tasks    []T
		executor func(t T) (R, error)
		opts     []DoOpt
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Do(tt.args.ctx, tt.args.tasks, tt.args.executor, tt.args.opts...)
			tt.requireErr(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}

func fibonacci(n int) (int, error) {
	if n < 0 {
		return 0, errors.New("negative")
	}
	if n < 2 {
		return 1, nil
	}
	r1, err := fibonacci(n - 1)
	if err != nil {
		return 0, err
	}

	r2, err := fibonacci(n - 2)
	if err != nil {
		return 0, err
	}

	return r1 + r2, nil
}
