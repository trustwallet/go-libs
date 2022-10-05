package metrics

import (
	"github.com/prometheus/client_golang/prometheus"

	"github.com/trustwallet/go-libs/logging"
)

func Register(labels prometheus.Labels, reg prometheus.Registerer, collectors ...prometheus.Collector) {
	registerer := prometheus.WrapRegistererWith(labels, reg)
	for _, c := range collectors {
		err := registerer.Register(c)
		if err != nil {
			if _, ok := err.(*prometheus.AlreadyRegisteredError); !ok {
				logging.GetLogger().WithError(err).
					Error("failed to register job duration metrics with prometheus")
			}
		}
	}
}
