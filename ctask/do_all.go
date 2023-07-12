package ctask

import (
	"context"
	"runtime"

	"golang.org/x/sync/errgroup"
)

type DoAllOpt func(cfg *DoAllConfig)
type DoAllConfig struct {
	WorkerNum int
}

type DoAllResp[R any] struct {
	Result R
	Error  error
}

// DoAll execute tasks using the given executor function for all the given tasks.
// It waits until all tasks are finished.
// The return value is a slice of Result or Error
//
// The max number of goroutines can optionally be specified using the option WithWorkerNum.
// By default, it is set to runtime.NumCPU()
func DoAll[Task any, Result any](
	ctx context.Context,
	tasks []Task,
	executor func(ctx context.Context, t Task) (Result, error),
	opts ...DoAllOpt,
) []DoAllResp[Result] {
	cfg := getDoAllConfigWithOptions(opts...)

	g, ctx := errgroup.WithContext(ctx)
	g.SetLimit(cfg.WorkerNum)
	results := make([]DoAllResp[Result], len(tasks))
	for idx, task := range tasks {
		idx, task := idx, task // retain current loop values to be used in goroutine
		g.Go(func() error {
			select {
			case <-ctx.Done():
				results[idx] = DoAllResp[Result]{Error: ctx.Err()}
				return nil
			default:
				res, err := executor(ctx, task)
				results[idx] = DoAllResp[Result]{
					Result: res,
					Error:  err,
				}
				return nil
			}
		})
	}
	if err := g.Wait(); err != nil {
		panic(err) // impossible to have error here
	}
	return results
}

func getDoAllConfigWithOptions(opts ...DoAllOpt) DoAllConfig {
	cfg := DoAllConfig{
		WorkerNum: runtime.NumCPU(),
	}
	for _, opt := range opts {
		opt(&cfg)
	}
	return cfg
}

func WithDoAllWorkerNum(num int) DoAllOpt {
	return func(cfg *DoAllConfig) {
		cfg.WorkerNum = num
	}
}
