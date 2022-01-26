package worker

import (
	"context"
	"sync"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/trustwallet/go-libs/logging"
)

const WorkerName = "simple_worker"

type SimpleWorker Instance

type Invoke func() error

type metricsWorker struct {
	interval time.Duration
	invoke   Invoke
	log      *logrus.Entry
}

func NewSimpleWoker(interval time.Duration, invoke Invoke) SimpleWorker {
	return &metricsWorker{
		interval: interval,
		invoke:   invoke,
		log:      logging.GetLogger().WithField("worker", WorkerName),
	}
}

func (m *metricsWorker) Start(ctx context.Context, wg *sync.WaitGroup) {
	New(WorkerName, func() {
		err := m.invoke()
		if err != nil {
			logging.GetLogger().WithField("worker", WorkerName).
				WithError(err).Error("error while invoking worker func")
		}
	}, m.interval, true).StartWithTicker(ctx, wg)
}

func (t *metricsWorker) Name() string {
	return WorkerName
}
