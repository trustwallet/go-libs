package worker

import (
	"context"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/trustwallet/go-libs/metrics"
)

type Builder interface {
	WithOptions(options *WorkerOptions) Builder
	WithStop(func() error) Builder
	Build() Worker
}

type builder struct {
	worker *worker
}

func NewWorkerBuilder(name string, workerFn func() error) Builder {
	return &builder{
		worker: &worker{
			name:     name,
			workerFn: workerFn,
			options:  DefaultWorkerOptions(1 * time.Minute),
		},
	}
}

func (b *builder) WithOptions(options *WorkerOptions) Builder {
	b.worker.options = options
	return b
}

func (b *builder) WithStop(stopFn func() error) Builder {
	b.worker.stopFn = stopFn
	return b
}

func (b *builder) Build() Worker {
	return b.worker
}

// Worker interface can be constructed using worker.NewBuilder("worker_name", workerFn).Build()
// or allows custom implementation (e.g. one-off jobs)
type Worker interface {
	Name() string
	Start(ctx context.Context, wg *sync.WaitGroup)
}

type worker struct {
	name     string
	workerFn func() error
	stopFn   func() error
	options  *WorkerOptions
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
				if w.stopFn != nil {
					log.WithField("worker", w.name).Info("stopping...")
					if err := w.stopFn(); err != nil {
						log.WithField("worker", w.name).WithError(err).Warn("error ocurred while stopping the worker")
					}
				}
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
