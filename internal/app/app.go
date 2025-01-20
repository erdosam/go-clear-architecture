package app

import (
	"fmt"
	"github.com/arendi-project/ba-version-2/config"
	v1 "github.com/arendi-project/ba-version-2/internal/controller/http/v1"
	"github.com/arendi-project/ba-version-2/pkg/httpserver"
	"github.com/arendi-project/ba-version-2/pkg/logger"
	"github.com/gin-gonic/gin"
	"os"
	"os/signal"
	"syscall"
)

var (
	logService logger.Interface
)

func Run(cfg *config.Config) {
	logService = newLogger()
	db := newRepo()
	defer db.Disconnect()

	listenAndServe(cfg)
}

func listenAndServe(cfg *config.Config) {
	// HTTP Server
	handler := gin.New()
	f := &v1.Feature{
		Carting: newCartingUseCase(),
		Order:   newOrderUseCase(),
		User:    newUserUseCase(),
	}
	m := &v1.Middleware{
		Authentication: newAuthenticationMiddleware(),
		Authorization:  newAuthorizationMiddleware(),
	}
	v1.InitRouter(handler, logService, f, m)
	hs := httpserver.New(handler, httpserver.Port(cfg.HTTP.Port))
	waitSignal(hs)
}

func waitSignal(httpServer *httpserver.Server) {
	var err error
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
