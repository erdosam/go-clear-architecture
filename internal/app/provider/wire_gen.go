// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package provider

import (
	"fmt"
	"github.com/arendi-project/ba-version-2/config"
	"github.com/arendi-project/ba-version-2/internal/controller/http/middleware"
	"github.com/arendi-project/ba-version-2/internal/controller/http/v1"
	"github.com/arendi-project/ba-version-2/internal/usecase"
	"github.com/arendi-project/ba-version-2/internal/usecase/dao"
	"github.com/arendi-project/ba-version-2/pkg/httpserver"
	"github.com/arendi-project/ba-version-2/pkg/logger"
	"github.com/arendi-project/ba-version-2/pkg/postgres"
	"github.com/go-playground/validator/v10"
	"github.com/google/wire"
	"log"
	"sync"
)

// Injectors from wire.go:

func NewRepo() *postgres.Postgres {
	config := provideSingletonConfig()
	loggerInterface := NewLogger()
	postgresPostgres := provideSingletonRepository(config, loggerInterface)
	return postgresPostgres
}

func NewLogger() logger.Interface {
	config := provideSingletonConfig()
	loggerInterface := provideSingletonLogger(config)
	return loggerInterface
}

func NewHttpServer() *httpserver.Server {
	config := provideSingletonConfig()
	loggerInterface := provideSingletonLogger(config)
	carting := newCartingUseCase()
	order := newOrderUseCase()
	user := newUserUseCase()
	feature := provideFeatures(carting, order, user)
	authentication := newAuthenticationMiddleware()
	authorization := newAuthorizationMiddleware()
	middleware := provideMiddlewares(authentication, authorization)
	handler := v1.NewRouterHandler(loggerInterface, feature, middleware)
	v := provideServerOptions(config)
	server := httpserver.New(handler, v...)
	return server
}

func newAuthenticationMiddleware() middleware.Authentication {
	user := newUserUseCase()
	config := provideSingletonConfig()
	loggerInterface := provideSingletonLogger(config)
	juno := provideJunoConfig(config)
	string2 := juno.ClientKeyFile
	authentication := middleware.NewJwtAuthentication(user, loggerInterface, string2)
	return authentication
}

func newAuthorizationMiddleware() middleware.Authorization {
	config := provideSingletonConfig()
	loggerInterface := provideSingletonLogger(config)
	authorization := middleware.NewAbacAuthorization(loggerInterface)
	return authorization
}

func newCartingUseCase() usecase.Carting {
	config := provideSingletonConfig()
	loggerInterface := provideSingletonLogger(config)
	postgresPostgres := provideSingletonRepository(config, loggerInterface)
	cartingDAO := newSingletonCartingDAO(loggerInterface, postgresPostgres)
	validate := provideValidator()
	carting := usecase.NewCartingUseCase(loggerInterface, cartingDAO, validate)
	return carting
}

func newOrderUseCase() usecase.Order {
	config := provideSingletonConfig()
	loggerInterface := provideSingletonLogger(config)
	postgresPostgres := provideSingletonRepository(config, loggerInterface)
	orderDAO := newSingletonOrderDAO(loggerInterface, postgresPostgres)
	cartingDAO := newSingletonCartingDAO(loggerInterface, postgresPostgres)
	order := usecase.NewOrderUseCase(loggerInterface, orderDAO, cartingDAO)
	return order
}

func newUserUseCase() usecase.User {
	config := provideSingletonConfig()
	loggerInterface := provideSingletonLogger(config)
	postgresPostgres := provideSingletonRepository(config, loggerInterface)
	userDAO := newSingletonUserDAO(loggerInterface, postgresPostgres)
	user := usecase.NewUserUseCase(loggerInterface, userDAO)
	return user
}

// wire.go:

type singletons struct {
	config     *config.Config
	log        logger.Interface
	db         *postgres.Postgres
	cartingDAO dao.CartingDAO
	orderDAO   dao.OrderDAO
	userDAO    dao.UserDAO
	once       struct {
		config     sync.Once
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
	commonSet = wire.NewSet(
		provideSingletonConfig,
		provideSingletonLogger,
		provideSingletonRepository,
		provideValidator,
	)
	daoSet = wire.NewSet(
		newSingletonCartingDAO,
		newSingletonOrderDAO,
		newSingletonUserDAO,
	)
	middlewareSet = wire.NewSet(
		newAuthenticationMiddleware,
		newAuthorizationMiddleware,
		provideMiddlewares,
	)
	featureSet = wire.NewSet(
		newCartingUseCase,
		newOrderUseCase,
		newUserUseCase,
		provideFeatures,
	)
)

func provideFeatures(c usecase.Carting, o usecase.Order, u usecase.User) *v1.Feature {
	return &v1.Feature{
		Carting: c,
		Order:   o,
		User:    u,
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

func provideSingletonConfig() *config.Config {
	s.once.config.Do(func() {
		conf, err := config.NewConfig()
		if err != nil {
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

func provideJunoConfig(cfg *config.Config) config.Juno {
	return cfg.Juno
}

func provideServerOptions(cfg *config.Config) []httpserver.Option {
	return []httpserver.Option{httpserver.Port(cfg.HTTP.Port)}
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
