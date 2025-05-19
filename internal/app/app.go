package app

import (
	"fmt"
	"github.com/erdosam/go-clear-architecture/internal/app/provider"
	"github.com/erdosam/go-clear-architecture/pkg/httpserver"
	"github.com/erdosam/go-clear-architecture/pkg/messaging"
	"github.com/erdosam/go-clear-architecture/pkg/storage"
	"os"
	"os/signal"
	"syscall"
)

func Run() {
	db := provider.NewRepo()
	defer db.Disconnect()

	subscribe()
	listenAndServe()
}

func subscribe() {
	ps := provider.NewPubsubSubscriber()
	ps.Init()
}

func listenAndServe() {
	// HTTP Server
	hs := provider.NewHttpServer()
	cs := provider.NewCloudStorage()
	cm := provider.NewCloudMessaging()
	waitSignal(hs, cs, cm)
}

func waitSignal(httpServer *httpserver.Server, storage storage.Storage, pubsub messaging.Pubsub) {
	var err error
	logService := provider.NewLogger()
	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM)
	// don't work on dev (run with "go run ./cmd/app"
	select {
	case s := <-interrupt:
		logService.Info("app - Run - signal: " + s.String())
	case err = <-httpServer.Notify():
		logService.Error(fmt.Errorf("app - Run - httpServer.Notify: %w", err))
	}
	// Shutdown
	err = httpServer.Shutdown()
	if err != nil {
		logService.Error(fmt.Errorf("app - Run - httpServer.Shutdown: %w", err))
	}
	_ = storage.Close()
	_ = pubsub.Close()
	//TODO if you need rpc server see https://github.com/evrone/go-clean-template
}
