package httplib

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"sync"

	log "github.com/sirupsen/logrus"
)

type Server struct {
	srv  *http.Server
	port int
}

// Run runs a server.
func (s Server) Run(wg *sync.WaitGroup) {
	wg.Add(1)

	go func() {
		log.WithField("port", s.port).Infof("Running HTTP server")

		if err := s.srv.ListenAndServe(); err != nil && errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("listen: %v", err)
		}

		wg.Done()
	}()
}

// Shutdown stops a server.
func (s Server) Shutdown(ctx context.Context) error {
	log.Info("Stopping HTTP server")

	return s.srv.Shutdown(ctx)
}

// NewHTTPServer returns an initialized HTTP server.
func NewHTTPServer(router http.Handler, port int) *Server {
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: router,
	}

	return &Server{
		srv:  srv,
		port: port,
	}
}
