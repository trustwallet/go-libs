package worker

import (
	"context"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/trustwallet/go-libs/metrics"
)

type Worker interface {
	Name() string
	Start(ctx context.Context, wg *sync.WaitGroup)
}

type worker struct {
	name     string
	workerFn func() error
	options  *WorkerOptions
}

func InitWorker(name string, options *WorkerOptions, workerFn func() error) Worker {
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
			log.WithField("worker", w.name).Info("run immediately")
			w.invoke()
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
				w.invoke()

				if w.options.RunConsequently {
					ticker = time.NewTicker(w.options.Interval)
				}
			}
		}
	}()
}

func (w *worker) invoke() {
	metric := w.options.PerformanceMetric
	if metric == nil {
		metric = &metrics.NullablePerformanceMetric{}
	}

	defer metric.Duration(metric.Start())
	err := w.workerFn()

	if err != nil {
		metric.Failure()
		log.WithField("worker", w.name).Error(err)
	} else {
		metric.Success()
	}
}
