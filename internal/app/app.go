package app

import (
	"fmt"
	"github.com/arendi-project/ba-version-2/config"
	"github.com/arendi-project/ba-version-2/internal/controller/http/middleware"
	v1 "github.com/arendi-project/ba-version-2/internal/controller/http/v1"
	"github.com/arendi-project/ba-version-2/internal/usecase"
	"github.com/arendi-project/ba-version-2/internal/usecase/dao"
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

var (
	comp       *mainComponent
	cartingDAO dao.CartingDAO
	orderDAO   dao.OrderDAO
	userDAO    dao.UserDAO
	httpServer *httpserver.Server
)

func Run(cfg *config.Config) {
	initAppComponents(cfg)
	defer comp.db.Disconnect()

	initDAOs()
	listenAndServe(cfg)
	waitSignal()
}

func initAppComponents(cfg *config.Config) {
	l := logger.New(cfg.Log.Level)
	// Repository
	pg, err := postgres.New(cfg.PG.URL)
	if err != nil {
		l.Fatal(fmt.Errorf("app.Run: %w", err))
	}
	comp = &mainComponent{log: l, db: pg}
}

func initDAOs() {
	cartingDAO = dao.NewCartingDAO(comp.log, comp.db)
	orderDAO = dao.NewOrderDAO(comp.log, comp.db)
	userDAO = dao.NewUserDAO(comp.log, comp.db)
}

func listenAndServe(cfg *config.Config) {
	// HTTP Server
	handler := gin.New()
	f := &v1.Feature{
		Carting: createCartingUseCase(),
		Order:   createOrderUseCase(),
		User:    createUserUseCase(),
	}
	m := &v1.Middleware{
		Authentication: createAuthenticationMiddleware(f.User, cfg.Juno.ClientKeyFile),
		Authorization:  createAuthorizationMiddleware(),
	}
	v1.NewRouter(handler, comp.log, f, m)
	httpServer = httpserver.New(handler, httpserver.Port(cfg.HTTP.Port))
}

func createCartingUseCase() usecase.Carting {
	return usecase.NewCartingUseCase(comp.log, cartingDAO)
}

func createOrderUseCase() usecase.Order {
	return usecase.NewOrderUseCase(comp.log, orderDAO, cartingDAO)
}

func createUserUseCase() usecase.User {
	return usecase.NewUserUseCase(comp.log, userDAO)
}

func createAuthenticationMiddleware(u usecase.User, f string) middleware.Authentication {
	return middleware.NewJwtAuthentication(u, comp.log, f)
}

func createAuthorizationMiddleware() middleware.Authorization {
	return middleware.NewAbacAuthorization(comp.log)
}

func waitSignal() {
	var err error
	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM)
	select {
	case s := <-interrupt:
		comp.log.Info("app - Run - signal: " + s.String())
	case err = <-httpServer.Notify():
		comp.log.Error(fmt.Errorf("app - Run - httpServer.Notify: %w", err))
	}
	// Shutdown
	err = httpServer.Shutdown()
	if err != nil {
		comp.log.Error(fmt.Errorf("app - Run - httpServer.Shutdown: %w", err))
	}
	//TODO if you need rpc server see https://github.com/evrone/go-clean-template
}
