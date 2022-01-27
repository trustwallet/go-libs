package worker_test

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/trustwallet/go-libs/worker"
	"gotest.tools/assert"
)

func TestWorkerWithDefaultOptions(t *testing.T) {
	counter := 0
	options := worker.DefaultWorkerOptions(100 * time.Millisecond)
	worker := worker.InitWorker("test", options, func() error {
		counter++
		return nil
	})
	wg := &sync.WaitGroup{}
	cxt := context.Background()

	worker.Start(cxt, wg)

	time.Sleep(350 * time.Millisecond)

	assert.Equal(t, 4, counter, "Should execute 3 times - 1st immidietly, and 3 after")
}

func TestWorkerStartsConsequently(t *testing.T) {
	counter := 0
	options := worker.DefaultWorkerOptions(100 * time.Millisecond).
		ShouldFinishBeforeNextStart()

	worker := worker.InitWorker("test", options, func() error {
		counter++
		time.Sleep(100 * time.Millisecond)
		return nil
	})
	wg := &sync.WaitGroup{}
	cxt := context.Background()

	worker.Start(cxt, wg)

	time.Sleep(350 * time.Millisecond)

	assert.Equal(t, 3, counter, "Should execute 2 times - 1st immidietly, and 2 after with delat 2sec between runs")
}
