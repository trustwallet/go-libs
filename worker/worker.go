package worker

import (
	"context"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

type Worker struct {
	workerFn       func()
	interval       time.Duration
	name           string
	runImmediately bool
}

func New(name string, workerFn func(), interval time.Duration, runImmediately bool) *Worker {
	return &Worker{
		name:           name,
		workerFn:       workerFn,
		interval:       interval,
		runImmediately: runImmediately,
	}
}

// StartConsequently waits for w.interval before each iteration
func (w *Worker) StartConsequently(ctx context.Context, wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		defer wg.Done()

		if w.runImmediately {
			w.workerFn()
		}

		for {
			select {
			case <-ctx.Done():
				log.WithField("service", w.name).Info("Stopped")
				return
			case <-time.After(w.interval):
				log.WithField("service", w.name).Info("Processing")
				w.workerFn()
			}
		}
	}()
}

// StartWithTicker executes the function with the provided interval
// In case execution takes longer than interval, next iteration start immediately
func (w *Worker) StartWithTicker(ctx context.Context, wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		defer wg.Done()

		ticker := time.NewTicker(w.interval)
		defer ticker.Stop()

		if w.runImmediately {
			w.workerFn()
		}

		for {
			select {
			case <-ctx.Done():
				log.WithField("service", w.name).Info("Stopped")
				return
			case <-ticker.C:
				log.WithField("service", w.name).Info("Processing")
				w.workerFn()
			}
		}
	}()
}
