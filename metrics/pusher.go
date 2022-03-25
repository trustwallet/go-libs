package metrics

import (
	"os"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/push"

	"github.com/trustwallet/go-libs/client"
)

type Pusher interface {
	Push() error
}

type pusher struct {
	pusher *push.Pusher
}

func NewPusher(pushgatewayURL, jobName string) Pusher {
	return &pusher{
		pusher: push.New(pushgatewayURL, jobName).
			Grouping("instance", instanceID()).
			Gatherer(prometheus.DefaultGatherer),
	}
}

func NewPusherWithCustomClient(pushgatewayURL, jobName string, client client.HTTPClient) Pusher {
	return &pusher{
		pusher: push.New(pushgatewayURL, string(jobName)).
			Grouping("instance", instanceID()).
			Gatherer(prometheus.DefaultGatherer).
			Client(client),
	}
}

func (p *pusher) Push() error {
	return p.pusher.Push()
}

func instanceID() string {
	instance := os.Getenv("DYNO")
	if instance == "" {
		instance = os.Getenv("INSTANCE_ID")
	}
	if instance == "" {
		instance = "local"
	}
	return instance
}
