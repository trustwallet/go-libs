package worker

import (
	"context"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

type worker struct {
	workerFn func()
	interval time.Duration
	name     string
}

func New(name string, workerFn func(), interval time.Duration) *worker {
	return &worker{
		name:     name,
		workerFn: workerFn,
		interval: interval,
	}
}

func (w *worker) Start(ctx context.Context, wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		defer wg.Done()

		for {
			select {
			case <-ctx.Done():
				log.WithField("service", w.name).Info("Stopped")
				return
			default:
				w.workerFn()
				time.Sleep(w.interval)
			}
		}
	}()
}
