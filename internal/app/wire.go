//go:build wireinject

package app

import (
	"fmt"
	"github.com/arendi-project/ba-version-2/config"
	"github.com/arendi-project/ba-version-2/internal/controller/http/middleware"
	"github.com/arendi-project/ba-version-2/internal/usecase"
	"github.com/arendi-project/ba-version-2/internal/usecase/dao"
	"github.com/arendi-project/ba-version-2/pkg/logger"
	"github.com/arendi-project/ba-version-2/pkg/postgres"
	"github.com/google/wire"
	"log"
	"sync"
)

type singletons struct {
	log        logger.Interface
	db         *postgres.Postgres
	cartingDAO dao.CartingDAO
	orderDAO   dao.OrderDAO
	userDAO    dao.UserDAO
	once       struct {
		log        sync.Once
		db         sync.Once
		cartingDAO sync.Once
		orderDAO   sync.Once
		userDAO    sync.Once
	}
}

var (
	s singletons
	// component dependencies
	commonSet = wire.NewSet(provideConfig, provideSingletonLogger, provideSingletonRepository)
	daoSet    = wire.NewSet(
		newSingletonCartingDAO,
		newSingletonOrderDAO,
		newSingletonUserDAO,
	)
)

func newRepo() *postgres.Postgres {
	wire.Build(provideConfig, newLogger, provideSingletonRepository)
	return nil
}

func newLogger() logger.Interface {
	wire.Build(provideConfig, provideSingletonLogger)
	return nil
}

func newAuthenticationMiddleware() middleware.Authentication {
	wire.Build(
		middleware.NewJwtAuthentication,
		newUserUseCase,
		provideJunoConfig,
		wire.FieldsOf(new(config.Juno), "ClientKeyFile"),
		commonSet,
	)
	return nil
}

func newAuthorizationMiddleware() middleware.Authorization {
	wire.Build(middleware.NewAbacAuthorization, commonSet)
	return nil
}

func newCartingUseCase() usecase.Carting {
	wire.Build(usecase.NewCartingUseCase, commonSet, daoSet)
	return nil
}

func newOrderUseCase() usecase.Order {
	wire.Build(usecase.NewOrderUseCase, commonSet, daoSet)
	return nil
}

func newUserUseCase() usecase.User {
	wire.Build(usecase.NewUserUseCase, commonSet, daoSet)
	return nil
}

func provideConfig() *config.Config {
	conf, err := config.NewConfig()
	if err != nil {
		// it uses "log" module since the logger is depends on config, avoid circular dependency
		log.Fatalf("Config error: %s", err)
	}
	return conf
}

func provideSingletonLogger(cfg *config.Config) logger.Interface {
	s.once.log.Do(func() {
		s.log = logger.New(cfg.Log.Level)
	})
	return s.log
}

func provideSingletonRepository(cfg *config.Config, l logger.Interface) *postgres.Postgres {
	s.once.db.Do(func() {
		pg, err := postgres.New(cfg.PG.URL)
		if err != nil {
			l.Fatal(fmt.Errorf("app.Run: %w", err))
		}
		s.db = pg
	})
	return s.db
}

func provideJunoConfig(cfg *config.Config) config.Juno {
	return cfg.Juno
}

func newSingletonCartingDAO(l logger.Interface, pg *postgres.Postgres) dao.CartingDAO {
	s.once.cartingDAO.Do(func() {
		s.cartingDAO = dao.NewCartingDAO(l, pg)
	})
	return s.cartingDAO
}

func newSingletonOrderDAO(l logger.Interface, pg *postgres.Postgres) dao.OrderDAO {
	s.once.orderDAO.Do(func() {
		s.orderDAO = dao.NewOrderDAO(l, pg)
	})
	return s.orderDAO
}

func newSingletonUserDAO(l logger.Interface, pg *postgres.Postgres) dao.UserDAO {
	s.once.userDAO.Do(func() {
		s.userDAO = dao.NewUserDAO(l, pg)
	})
	return s.userDAO
}
