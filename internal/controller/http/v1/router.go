package v1

import (
	"net/http"

	"github.com/erdosam/go-clear-architecture/internal/controller/http/middleware"
	"github.com/erdosam/go-clear-architecture/internal/entity"
	"github.com/erdosam/go-clear-architecture/internal/usecase"
	"github.com/erdosam/go-clear-architecture/pkg/logger"
	"github.com/gin-gonic/gin"
)

type Feature struct {
	Ping usecase.Ping
	User usecase.User
}

type Middleware struct {
	Authentication middleware.Authentication
	Authorization  middleware.Authorization
}

type routerHandler struct {
	*gin.Engine
	log     logger.Interface
	feature *Feature
	access  middleware.Authorization
}

func NewRouterHandler(log logger.Interface, f *Feature, m *Middleware) http.Handler {
	handler := &routerHandler{
		gin.New(),
		log,
		f,
		m.Authorization,
	}
	handler.Use(gin.Logger())
	handler.Use(gin.Recovery())
	// Swagger
	// TODO if you need this see https://github.com/evrone/go-clean-template

	// K8s probe
	handler.GET("/ping", func(c *gin.Context) { c.JSON(http.StatusOK, "pong") })
	handler.HEAD("/ping", func(c *gin.Context) { c.Status(http.StatusOK) })
	handler.PATCH("/load-policy", m.Authorization.LoadPolicy)

	// Prometheus metrics
	// TODO if you need this see https://github.com/evrone/go-clean-template

	v1 := handler.Group("/v1", m.Authentication.Authenticate)
	{
		handler.initPingRoutes(v1)
		//TODO define other routes here
	}
	return handler
}

func getIdentity(c *gin.Context) entity.User {
	user, _ := c.Get(middleware.IdentityContextKey)
	return user.(entity.User)
}

func setAccessResource(c *gin.Context, value any) {
	c.Set(middleware.ResourceContextKey, value)
}
