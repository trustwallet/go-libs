package gin

import (
	"context"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// SetupGracefulShutdown blocks execution until interruption command sent
func SetupGracefulShutdown(ctx context.Context, port string, engine *gin.Engine) {
	server := &http.Server{
		Addr:    ":" + port,
		Handler: engine,
	}

	defer func() {
		if err := server.Shutdown(ctx); err != nil {
			log.Info("Server Shutdown: ", err)
		}
	}()

	signalForExit := make(chan os.Signal, 1)
	signal.Notify(signalForExit,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	go func() {
		switch err := server.ListenAndServe(); err {
		case http.ErrServerClosed:
			log.Info("server closed")
		default:
			log.Error("Application failed ", err)
		}
	}()
	log.WithFields(log.Fields{"bind": port}).Info("Running application")

	stop := <-signalForExit
	log.Info("Stop signal Received ", stop)
	log.Info("Waiting for all jobs to stop")
}

// SetupGracefulServeWithUnixFile blocks execution until interruption command sent
func SetupGracefulServeWithUnixFile(ctx context.Context, engine *gin.Engine, unixFile string) {
	_, err := os.Create("/tmp/app-initialized")
	if err != nil {
		log.WithError(err).Error("failed to create file /tmp/app-initialized")
		return
	}

	defer func() {
		if err != nil {
			log.Error(err)
		}
	}()

	listener, err := net.Listen("unix", unixFile)
	if err != nil {
		return
	}

	defer func() { _ = listener.Close() }()
	defer func() { _ = os.Remove(unixFile) }()

	server := &http.Server{
		Handler: engine,
	}

	defer func() {
		if err := server.Shutdown(ctx); err != nil {
			log.Info("Server Shutdown: ", err)
		}
	}()

	signalForExit := make(chan os.Signal, 1)
	signal.Notify(signalForExit,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	go func() {
		log.Debugf("Listening and serving HTTP on unix:/%s", unixFile)
		err = server.Serve(listener)
	}()

	stop := <-signalForExit
	log.Info("Stop signal Received ", stop)
}
