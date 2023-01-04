package ctask

import (
	"context"
	"runtime"

	"golang.org/x/sync/errgroup"
)

type DoOpt func(cfg *DoConfig)
type DoConfig struct {
	WorkerNum int
}

// Do execute tasks using the given executor function,
// and return the results in the same order as the given tasks respectively.
// It stops executing remaining tasks after any first error is encountered.
//
// The max number of goroutines can optionally be specified using the option WithWorkerNum.
// By default, it is set to runtime.NumCPU()
func Do[Task any, Result any](
	ctx context.Context,
	tasks []Task,
	executor func(ctx context.Context, t Task) (Result, error),
	opts ...DoOpt,
) ([]Result, error) {
	cfg := getConfigWithOptions(opts...)

	g, ctx := errgroup.WithContext(ctx)
	g.SetLimit(cfg.WorkerNum)
	results := make([]Result, len(tasks))
	for idx, task := range tasks {
		idx, task := idx, task // retain current loop values to be used in goroutine
		g.Go(func() error {
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
				res, err := executor(ctx, task)
				if err != nil {
					return err
				}
				results[idx] = res
				return nil
			}
		})
	}
	if err := g.Wait(); err != nil {
		return nil, err
	}
	return results, nil
}

func getConfigWithOptions(opts ...DoOpt) DoConfig {
	cfg := DoConfig{
		WorkerNum: runtime.NumCPU(),
	}
	for _, opt := range opts {
		opt(&cfg)
	}
	return cfg
}

func WithWorkerNum(num int) DoOpt {
	return func(cfg *DoConfig) {
		cfg.WorkerNum = num
	}
}
