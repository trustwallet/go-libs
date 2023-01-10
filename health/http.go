package health

import (
	"context"
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"
)

const (
	defaultHealthCheckRoute    = "/health"
	defaultReadinessCheckRoute = "/ready"
	defaultPort                = 4444
)

type CheckFunc func() error

type Option func(*server)

type server struct {
	healthCheckRoute    string
	readinessCheckRoute string
	port                int
	healthCheckFunc     CheckFunc
	readinessCheckFunc  CheckFunc
}

func WithHealthCheckRoute(route string) Option {
	return func(s *server) {
		s.healthCheckRoute = route
	}
}

func WithReadinessCheckRoute(route string) Option {
	return func(s *server) {
		s.readinessCheckRoute = route
	}
}

func WithPort(port int) Option {
	return func(s *server) {
		s.port = port
	}
}

func WithHealthCheckFunc(healthCheckFunc CheckFunc) Option {
	return func(s *server) {
		s.healthCheckFunc = healthCheckFunc
	}
}

func WithReadinessCheckFunc(readinessCheckFunc CheckFunc) Option {
	return func(s *server) {
		s.readinessCheckFunc = readinessCheckFunc
	}
}

func handle(handler *http.ServeMux, route string, handleFunc CheckFunc) {
	handler.HandleFunc(route, func(w http.ResponseWriter, r *http.Request) {
		if handleFunc == nil {
			w.WriteHeader(http.StatusOK)
			return
		}

		if err := handleFunc(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	})
}

// StartHealthCheckServer starts a HTTP server to handle health check and readiness check requests.
func StartHealthCheckServer(ctx context.Context, opts ...Option) error {
	hcServer := &server{
		healthCheckRoute:    defaultHealthCheckRoute,
		readinessCheckRoute: defaultReadinessCheckRoute,
		port:                defaultPort,
	}

	for _, opt := range opts {
		opt(hcServer)
	}

	handler := http.NewServeMux()
	handle(handler, hcServer.healthCheckRoute, hcServer.healthCheckFunc)
	handle(handler, hcServer.readinessCheckRoute, hcServer.readinessCheckFunc)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", hcServer.port),
		Handler: handler,
	}

	go func() {
		<-ctx.Done()
		if err := srv.Shutdown(ctx); err != nil {
			log.Info("server shutdown: ", err)
		}
	}()

	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		return err
	}

	return nil
}
