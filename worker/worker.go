package worker

import (
	"context"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

type Worker interface {
	Start(ctx context.Context, wg *sync.WaitGroup)
	Name() string
}

type worker struct {
	name     string
	workerFn func() error
	options  WorkerOptions
}

func InitWorker(name string, options WorkerOptions, workerFn func() error) Worker {
	return &worker{
		name:     name,
		options:  options,
		workerFn: workerFn,
	}
}

func (w *worker) Name() string {
	return w.name
}

func (w *worker) Start(ctx context.Context, wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		defer wg.Done()

		ticker := time.NewTicker(w.options.Interval)
		defer ticker.Stop()

		if w.options.RunImmediately {
			w.workerFn()
		}

		for {
			select {
			case <-ctx.Done():
				log.WithField("worker", w.name).Info("stopped")
				return
			case <-ticker.C:
				if w.options.RunConsequently {
					ticker.Stop()
				}

				log.WithField("worker", w.name).Info("processing")
				w.workerFn()

				if w.options.RunConsequently {
					ticker = time.NewTicker(w.options.Interval)
				}
			}
		}
	}()
}

// StartConsequently waits for w.interval before each iteration
// Deprecated: User Start() method
func (w *worker) StartConsequently(ctx context.Context, wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		defer wg.Done()

		if w.options.RunImmediately {
			w.workerFn()
		}

		for {
			select {
			case <-ctx.Done():
				log.WithField("worker", w.name).Info("stopped")
				return
			case <-time.After(w.options.Interval):
				log.WithField("worker", w.name).Info("processing")
				w.workerFn()
			}
		}
	}()
}

// StartWithTicker executes the function with the provided interval
// In case execution takes longer than interval, next iteration start immediately
// Deprecated: User Start() method
func (w *worker) StartWithTicker(ctx context.Context, wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		defer wg.Done()

		ticker := time.NewTicker(w.options.Interval)
		defer ticker.Stop()

		if w.options.RunImmediately {
			w.workerFn()
		}

		for {
			select {
			case <-ctx.Done():
				log.WithField("worker", w.name).Info("stopped")
				return
			case <-ticker.C:
				log.WithField("worker", w.name).Info("processing")
				w.workerFn()
			}
		}
	}()
}
