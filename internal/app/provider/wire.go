//go:build wireinject

package provider

import (
	"fmt"
	"github.com/casbin/casbin/v2"
	"github.com/erdosam/go-clear-architecture/config"
	"github.com/erdosam/go-clear-architecture/internal/controller/http/middleware"
	v1 "github.com/erdosam/go-clear-architecture/internal/controller/http/v1"
	"github.com/erdosam/go-clear-architecture/internal/controller/pubsub"
	"github.com/erdosam/go-clear-architecture/internal/usecase"
	"github.com/erdosam/go-clear-architecture/internal/usecase/dao"
	"github.com/erdosam/go-clear-architecture/pkg/httpserver"
	"github.com/erdosam/go-clear-architecture/pkg/logger"
	"github.com/erdosam/go-clear-architecture/pkg/messaging"
	"github.com/erdosam/go-clear-architecture/pkg/postgres"
	"github.com/go-playground/validator/v10"
	"github.com/google/wire"
	"log"
	"sync"
)

type singletons struct {
	config   *config.Config
	log      logger.Interface
	db       *postgres.Postgres
	enforcer *casbin.Enforcer
	userDAO  dao.UserDAO
	pubsub   messaging.Pubsub
	once     struct {
		config   sync.Once
		log      sync.Once
		db       sync.Once
		userDAO  sync.Once
		enforcer sync.Once
		pubsub   sync.Once
	}
}

var (
	s singletons
	// component dependencies
	commonSet = wire.NewSet(
		provideSingletonConfig,
		provideSingletonLogger,
		provideSingletonRepository,
		provideValidator,
		provideSingletonPubsub,
	)
	daoSet = wire.NewSet(
		newSingletonUserDAO,
	)
	middlewareSet = wire.NewSet(
		newAuthenticationMiddleware,
		newAuthorizationMiddleware,
		provideMiddlewares,
	)
	featureSet = wire.NewSet(
		newPingUsecase,
		provideFeatures,
	)
)

func NewRepo() *postgres.Postgres {
	wire.Build(provideSingletonConfig, NewLogger, provideSingletonRepository)
	return nil
}

func NewLogger() logger.Interface {
	wire.Build(provideSingletonConfig, provideSingletonLogger)
	return nil
}

func NewHttpServer() *httpserver.Server {
	wire.Build(
		httpserver.New,
		v1.NewRouterHandler,
		provideServerOptions,
		commonSet,
		featureSet,
		middlewareSet,
	)
	return nil
}

func NewPubsubSubscriber() *pubsub.SubscriptionHandler {
	wire.Build(
		pubsub.NewSubscriptionsHandler,
		commonSet,
	)
	return nil
}

func provideFeatures(ping usecase.Ping) *v1.Feature {
	return &v1.Feature{
		Ping: ping,
	}
}

func provideMiddlewares(a1 middleware.Authentication, a2 middleware.Authorization) *v1.Middleware {
	return &v1.Middleware{
		Authentication: a1,
		Authorization:  a2,
	}
}

func provideValidator() *validator.Validate {
	return validator.New(validator.WithRequiredStructEnabled())
}

func newAuthenticationMiddleware() middleware.Authentication {
	wire.Build(
		middleware.NewJwtAuthentication,
		newUserUseCase,
		commonSet,
	)
	return nil
}

func newAuthorizationMiddleware() middleware.Authorization {
	wire.Build(middleware.NewAbacAuthorization, commonSet)
	return nil
}

func newUserUseCase() usecase.User {
	wire.Build(usecase.NewUserUseCase, commonSet, daoSet)
	return nil
}

func newPingUsecase() usecase.Ping {
	wire.Build(usecase.NewPingUsecase, commonSet)
	return nil
}

func provideSingletonConfig() *config.Config {
	s.once.config.Do(func() {
		conf, err := config.NewConfig()
		if err != nil {
			// it uses "log" module since the logger is depends on config, avoid circular dependency
			log.Fatalf("Config error: %s", err)
		}
		s.config = conf
	})
	return s.config
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

func provideSingletonPubsub(l logger.Interface, cfg *config.Config) messaging.Pubsub {
	s.once.pubsub.Do(func() {
		s.pubsub = messaging.NewPubsub(cfg, l)
	})
	return s.pubsub
}

func provideServerOptions(cfg *config.Config) []httpserver.Option {
	return []httpserver.Option{httpserver.Port(cfg.HTTP.Port)}
}

func newSingletonUserDAO(l logger.Interface, pg *postgres.Postgres) dao.UserDAO {
	s.once.userDAO.Do(func() {
		s.userDAO = dao.NewUserDAO(l, pg)
	})
	return s.userDAO
}
