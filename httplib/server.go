package httplib

import (
	"context"
	"errors"
	"net/http"
	"sync"

	log "github.com/sirupsen/logrus"
)

type Server interface {
	Run(ctx context.Context, wg *sync.WaitGroup)
}

type api struct {
	router http.Handler
	port   string
}

func NewHTTPServer(router http.Handler, port string) Server {
	return &api{
		router: router,
		port:   port,
	}
}

func (a *api) Run(ctx context.Context, wg *sync.WaitGroup) {
	a.serve(ctx, wg)
}

func (a *api) serve(ctx context.Context, wg *sync.WaitGroup) {
	wg.Add(1)

	server := &http.Server{
		Addr:    ":" + a.port,
		Handler: a.router,
	}

	serverStopped := make(chan struct{})

	go func() {
		if err := server.ListenAndServe(); err != nil && errors.Is(err, http.ErrServerClosed) {
			log.Error("Server ListenAndServe: ", err)
			serverStopped <- struct{}{}
		}
	}()

	log.WithFields(log.Fields{"bind": a.port}).Info("Starting the API server")

	go func() {
		defer func() { wg.Done() }()

		select {
		case <-ctx.Done():
			log.Info("Shutting down the server")

			if err := server.Shutdown(context.Background()); err != nil {
				log.Info("Server Shutdown: ", err)
			}

			return
		case <-serverStopped:
			return
		}
	}()
}
