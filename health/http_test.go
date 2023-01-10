package health_test

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	. "github.com/trustwallet/go-libs/health"
)

func TestStartHealthCheckServer(t *testing.T) {
	tests := []struct {
		name                string
		healthCheckFunc     func() error
		readinessCheckFunc  func() error
		healthCheckRoute    string
		readinessCheckRoute string
		port                int
		expHealthy          bool
		expReady            bool
	}{
		{
			name:       "default case",
			expHealthy: true,
			expReady:   true,
		},
		{
			name:               "not healthy",
			healthCheckFunc:    func() error { return errors.New("health check") },
			readinessCheckFunc: func() error { return nil },
			port:               1111,
			expHealthy:         false,
			expReady:           true,
		},
		{
			name:               "not ready",
			healthCheckFunc:    func() error { return nil },
			readinessCheckFunc: func() error { return errors.New("health check") },
			port:               2222,
			expHealthy:         true,
			expReady:           false,
		},
		{
			name:                "custom routes and port",
			healthCheckFunc:     func() error { return nil },
			readinessCheckFunc:  func() error { return nil },
			healthCheckRoute:    "/custom-health",
			readinessCheckRoute: "/custom-ready",
			port:                3333,
			expHealthy:          true,
			expReady:            true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			var opts []Option
			if test.healthCheckFunc != nil {
				opts = append(opts, WithHealthCheckFunc(test.healthCheckFunc))
			}

			if test.readinessCheckFunc != nil {
				opts = append(opts, WithReadinessCheckFunc(test.readinessCheckFunc))
			}

			if test.healthCheckRoute != "" {
				opts = append(opts, WithHealthCheckRoute(test.healthCheckRoute))
			}

			if test.readinessCheckRoute != "" {
				opts = append(opts, WithReadinessCheckRoute(test.readinessCheckRoute))
			}

			if test.port != 0 {
				opts = append(opts, WithPort(test.port))
			}

			port := 4444
			if test.port != 0 {
				port = test.port
			}

			healthRoute := "/health"
			if test.healthCheckRoute != "" {
				healthRoute = test.healthCheckRoute
			}

			healthURL := fmt.Sprintf("http://:%d/%s", port, healthRoute)

			readinessRoute := "/ready"
			if test.readinessCheckRoute != "" {
				readinessRoute = test.readinessCheckRoute
			}

			readinessURL := fmt.Sprintf("http://:%d/%s", port, readinessRoute)

			go func() {
				assert.NoError(t, StartHealthCheckServer(ctx, opts...))
			}()
			waitForServerToStart(t, healthURL, 20*time.Millisecond, 1*time.Second)

			resp, err := http.Get(healthURL)
			assert.NoError(t, err)
			assert.True(t, (test.expHealthy && resp.StatusCode == http.StatusOK) || (!test.expHealthy && resp.StatusCode != http.StatusOK))

			resp, err = http.Get(readinessURL)
			assert.NoError(t, err)
			assert.True(t, (test.expReady && resp.StatusCode == http.StatusOK) || (!test.expReady && resp.StatusCode != http.StatusOK))

			cancel()
		})
	}
}

func waitForServerToStart(t *testing.T, url string, interval time.Duration, timeout time.Duration) {
	tick := time.NewTicker(interval)
	defer tick.Stop()
	now := time.Now()
	for {
		if time.Since(now) > timeout {
			t.Fatal("timeout to connect to server")
			return
		}

		<-tick.C
		if _, err := http.Get(url); err == nil {
			return
		}
	}
}

func TestServerClosedOnContextCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		assert.NoError(t, StartHealthCheckServer(ctx))
	}()
	waitForServerToStart(t, "http://:4444/health", 20*time.Millisecond, 1*time.Second)

	cancel()
	time.Sleep(time.Millisecond * 100)
	_, err := http.Get("http://:4444/health")
	assert.Error(t, err) // server was shut down
}
