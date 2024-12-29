package app

import (
	"fmt"
	"github.com/arendi-project/ba-version-2/config"
	v1 "github.com/arendi-project/ba-version-2/internal/controller/http/v1"
	"github.com/arendi-project/ba-version-2/internal/usecase/dao"
	"github.com/arendi-project/ba-version-2/internal/usecase/form"
	"github.com/arendi-project/ba-version-2/pkg/httpserver"
	"github.com/arendi-project/ba-version-2/pkg/logger"
	"github.com/arendi-project/ba-version-2/pkg/postgres"
	"github.com/gin-gonic/gin"
	"os"
	"os/signal"
	"syscall"
)

type mainComponent struct {
	log logger.Interface
	db  *postgres.Postgres
}

var comp mainComponent

func Run(cfg *config.Config) {
	log := logger.New(cfg.Log.Level)

	// Repository
	pg, err := postgres.New(cfg.PG.URL)
	if err != nil {
		log.Fatal(fmt.Errorf("app.Run: %w", err))
	}
	defer pg.Disconnect()

	// HTTP Server
	comp = mainComponent{log: log, db: pg}
	handler := gin.New()
	feature := v1.Feature{
		Carting: getCartingUseCase(),
	}
	v1.NewRouter(handler, log, feature)
	httpServer := httpserver.New(handler, httpserver.Port(cfg.HTTP.Port))

	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM)
	select {
	case s := <-interrupt:
		log.Info("app - Run - signal: " + s.String())
	case err = <-httpServer.Notify():
		log.Error(fmt.Errorf("app - Run - httpServer.Notify: %w", err))
	}
	// Shutdown
	err = httpServer.Shutdown()
	if err != nil {
		log.Error(fmt.Errorf("app - Run - httpServer.Shutdown: %w", err))
	}
	//TODO if you need rmqServer see https://github.com/evrone/go-clean-template
}

func getCartingUseCase() *form.CartingUseCase {
	cd := dao.NewCartingDAO(comp.log, comp.db)
	return form.NewCarting(comp.log, cd)
}
