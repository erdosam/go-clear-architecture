//go:build wireinject

package provider

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
	pgadapter "github.com/casbin/casbin-pg-adapter"
	"github.com/casbin/casbin/v2"
	"github.com/go-playground/validator/v10"
	"github.com/google/wire"
	"log"
	"sync"
	"time"
)

type singletons struct {
	config      *config.Config
	log         logger.Interface
	db          *postgres.Postgres
	enforcer    *casbin.Enforcer
	categoryDAO dao.TrashCategoryDAO
	cartingDAO  dao.CartingDAO
	orderDAO    dao.OrderDAO
	userDAO     dao.UserDAO
	once        struct {
		config      sync.Once
		log         sync.Once
		db          sync.Once
		categoryDAO sync.Once
		cartingDAO  sync.Once
		orderDAO    sync.Once
		userDAO     sync.Once
		enforcer    sync.Once
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
		newSingletonTrashCategoryDAO,
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
		newCategoryUseCase,
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

func provideFeatures(cat usecase.Category, c usecase.Carting, o usecase.Order, u usecase.User) *v1.Feature {
	return &v1.Feature{
		Category: cat,
		Carting:  c,
		Order:    o,
		User:     u,
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
		provideJunoConfig,
		wire.FieldsOf(new(config.Juno), "ClientKeyFile"),
		commonSet,
	)
	return nil
}

func newAuthorizationMiddleware() middleware.Authorization {
	wire.Build(middleware.NewAbacAuthorization, commonSet, provideSingletonCasbinEnforcer)
	return nil
}

func provideSingletonCasbinEnforcer(cfg *config.Config, l logger.Interface) *casbin.Enforcer {
	s.once.enforcer.Do(func() {
		pga, err := pgadapter.NewAdapter(cfg.PG.URL)
		if err != nil {
			l.Fatal(err)
		}
		enforcer, err := casbin.NewEnforcer(cfg.Casbin.ModelFile, pga)
		if err != nil {
			l.Fatal(err)
		}
		err = enforcer.LoadPolicy()
		if err != nil {
			l.Fatal(err)
		}
		ticker := time.NewTicker(time.Duration(cfg.Casbin.LoadPolicyInterval) * time.Second)
		go func() {
			for {
				select {
				case <-ticker.C:
					l.Info("Load policy")
					err := enforcer.LoadPolicy()
					if err != nil {
						l.Error("Error load policy: %s", err.Error())
					}
				}
			}
		}()
		s.enforcer = enforcer
	})
	return s.enforcer
}

func newCategoryUseCase() usecase.Category {
	wire.Build(usecase.NewCategoryUsecase, commonSet, daoSet)
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

func provideJunoConfig(cfg *config.Config) config.Juno {
	return cfg.Juno
}

func provideServerOptions(cfg *config.Config) []httpserver.Option {
	return []httpserver.Option{httpserver.Port(cfg.HTTP.Port)}
}

func newSingletonTrashCategoryDAO(l logger.Interface, pg *postgres.Postgres) dao.TrashCategoryDAO {
	s.once.categoryDAO.Do(func() {
		s.categoryDAO = dao.NewTrashCategoryDAO(l, pg)
	})
	return s.categoryDAO
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
