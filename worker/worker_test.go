package worker_test

import (
	"context"
	"sync"
	"testing"
	"time"

	"gotest.tools/assert"

	"github.com/trustwallet/go-libs/worker"
)

func TestWorkerWithDefaultOptions(t *testing.T) {
	counter := 0
	worker := worker.NewWorkerBuilder("test", func() error {
		counter++
		return nil
	}).WithOptions(worker.DefaultWorkerOptions(100 * time.Millisecond)).Build()

	wg := &sync.WaitGroup{}
	ctx, cancel := context.WithTimeout(context.Background(), 350*time.Millisecond)
	defer cancel()

	worker.Start(ctx, wg)

	wg.Wait()

	assert.Equal(t, 4, counter, "Should execute 4 times - 1st immediately, and 3 after")
}

func TestWorkerStartsConsequently(t *testing.T) {
	counter := 0
	options := worker.DefaultWorkerOptions(100 * time.Millisecond)
	options.RunConsequently = true

	worker := worker.NewWorkerBuilder("test", func() error {
		time.Sleep(100 * time.Millisecond)
		counter++
		return nil
	}).WithOptions(options).Build()

	wg := &sync.WaitGroup{}
	ctx, cancel := context.WithTimeout(context.Background(), 350*time.Millisecond)
	defer cancel()

	worker.Start(ctx, wg)

	wg.Wait()

	assert.Equal(t, 3, counter, "Should execute 3 times - 1st immediately, and 2 after with delay between runs")
}

func TestWorkerStartsWithoutExecution(t *testing.T) {
	counter := 0
	options := worker.DefaultWorkerOptions(100 * time.Millisecond)
	options.Interval = -1

	worker := worker.NewWorkerBuilder("test", func() error {
		time.Sleep(100 * time.Millisecond)
		counter++
		return nil
	}).WithOptions(options).Build()

	wg := &sync.WaitGroup{}
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	worker.Start(ctx, wg)

	wg.Wait()

	assert.Equal(t, 0, counter, "Should never be executed")
}
