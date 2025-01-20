package app

import (
	"fmt"
	"github.com/arendi-project/ba-version-2/internal/app/provider"
	"github.com/arendi-project/ba-version-2/pkg/httpserver"
	"os"
	"os/signal"
	"syscall"
)

func Run() {
	db := provider.NewRepo()
	defer db.Disconnect()

	listenAndServe()
}

func listenAndServe() {
	// HTTP Server
	hs := provider.NewHttpServer()
	waitSignal(hs)
}

func waitSignal(httpServer *httpserver.Server) {
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
	//TODO if you need rpc server see https://github.com/evrone/go-clean-template
}
