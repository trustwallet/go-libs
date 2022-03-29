package metrics

import (
	"net/http"
	"os"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/push"

	"github.com/trustwallet/go-libs/client"
)

type MetricsPusherClient struct {
	client client.Request
}

func NewMetricsPusherClient(pushURL, key string, errorHandler client.HttpErrorHandler) *MetricsPusherClient {
	client := client.InitClient(pushURL, errorHandler)
	client.AddHeader("X-API-Key", key)

	return &MetricsPusherClient{
		client: client,
	}
}

func (c *MetricsPusherClient) Do(req *http.Request) (*http.Response, error) {
	for key, value := range c.client.Headers {
		req.Header.Set(key, value)
	}
	return c.client.HttpClient.Do(req)
}

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
