package ctask

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestDoAll(t *testing.T) {
	type T = int // task type
	type R = int // result type

	type args struct {
		ctx        context.Context
		ctxTimeout time.Duration
		tasks      []T
		executor   func(ctx context.Context, t T) (R, error)
		opts       []DoAllOpt
	}
	tests := []struct {
		name string
		args args
		want []DoAllResp[R]
	}{
		{
			name: "happy path",
			args: args{
				ctx:      context.Background(),
				tasks:    []T{0, 1, 2, 3, 4, 5, 6},
				executor: fibonacci,
				opts:     nil,
			},
			want: []DoAllResp[R]{
				{Result: 1},
				{Result: 1},
				{Result: 2},
				{Result: 3},
				{Result: 5},
				{Result: 8},
				{Result: 13},
			},
		},
		{
			name: "empty slice",
			args: args{
				ctx:      context.Background(),
				tasks:    nil,
				executor: fibonacci,
				opts:     nil,
			},
			want: []DoAllResp[R]{},
		},
		{
			name: "error path",
			args: args{
				ctx:      context.Background(),
				tasks:    []T{0, 1, 2, 1, -1, 5},
				executor: fibonacci,
				opts:     []DoAllOpt{WithDoAllWorkerNum(1)},
			},
			want: []DoAllResp[R]{
				{Result: 1},
				{Result: 1},
				{Result: 2},
				{Result: 1},
				{Error: errors.New("negative")},
				{Result: 8},
			},
		},
		{
			name: "slow functions should return context deadline exceeded error",
			args: args{
				ctx:        context.Background(),
				ctxTimeout: 20 * time.Millisecond,
				tasks: []T{
					10, 1000, 10, 5000, 10,
				},
				executor: func(ctx context.Context, t T) (R, error) {
					select {
					case <-ctx.Done():
						return 0, ctx.Err()
					case <-time.After(time.Duration(t) * time.Millisecond):
						return 1, nil
					}
				},
				opts: []DoAllOpt{WithDoAllWorkerNum(5)},
			},
			want: []DoAllResp[R]{
				{Result: 1}, {Error: context.DeadlineExceeded}, {Result: 1}, {Error: context.DeadlineExceeded}, {Result: 1},
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
				opts: []DoAllOpt{WithDoAllWorkerNum(20)},
			},
			want: []DoAllResp[R]{
				{Result: 1}, {Result: 1}, {Result: 1}, {Result: 1}, {Result: 1},
				{Result: 1}, {Result: 1}, {Result: 1}, {Result: 1}, {Result: 1},
				{Result: 1}, {Result: 1}, {Result: 1}, {Result: 1}, {Result: 1},
				{Result: 1}, {Result: 1}, {Result: 1}, {Result: 1}, {Result: 1},
				{Result: 1}, {Result: 1}, {Result: 1}, {Result: 1}, {Result: 1},
				{Result: 1}, {Result: 1}, {Result: 1}, {Result: 1}, {Result: 1},
				{Result: 1}, {Result: 1}, {Result: 1}, {Result: 1}, {Result: 1},
				{Result: 1}, {Result: 1}, {Result: 1}, {Result: 1}, {Result: 1},
				{Result: 1}, {Result: 1}, {Result: 1}, {Result: 1}, {Result: 1},
				{Result: 1}, {Result: 1}, {Result: 1}, {Result: 1}, {Result: 1},
			},
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
			got := DoAll(ctx, tt.args.tasks, tt.args.executor, tt.args.opts...)
			require.Equal(t, tt.want, got)
		})
	}
}
