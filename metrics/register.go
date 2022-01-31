package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/trustwallet/go-libs/logging"
)

func Register(collectors ...prometheus.Collector) {
	for _, c := range collectors {
		err := prometheus.Register(c)
		if err != nil &&
			err.Error() != "duplicate metrics collector registration attempted" {

			logging.GetLogger().WithError(err).
				Error("failed to register job duration metrics with prometheus")
		}
	}
}
