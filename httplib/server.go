package httplib

import (
	"context"
	"errors"
	"net/http"
	"sync"

	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"

	log "github.com/sirupsen/logrus"
)

type Server interface {
	Run(ctx context.Context, wg *sync.WaitGroup)
}

type api struct {
	router http.Handler
	port   string
	h2c    bool
}

func NewHTTPServer(router http.Handler, port string) Server {
	return &api{
		router: router,
		port:   port,
		h2c:    false,
	}
}

func NewH2CServer(router http.Handler, port string) Server {
	return &api{
		router: router,
		port:   port,
		h2c:    true,
	}
}

func (a *api) Run(ctx context.Context, wg *sync.WaitGroup) {
	a.serve(ctx, wg)
}

func (a *api) serve(ctx context.Context, wg *sync.WaitGroup) {
	wg.Add(1)

	h2s := &http2.Server{}
	h1d, h2d := a.router, h2c.NewHandler(a.router, h2s)

	server := &http.Server{
		Addr:    ":" + a.port,
		Handler: h1d,
	}
	if a.h2c {
		server.Handler = h2d
	}

	serverStopped := make(chan struct{})

	go func() {
		if err := server.ListenAndServe(); err != nil && errors.Is(err, http.ErrServerClosed) {
			log.WithError(err).Debug("Server ListenAndServe")
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
