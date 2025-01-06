package v1

import (
	"github.com/arendi-project/ba-version-2/internal/controller/http/middleware"
	"github.com/arendi-project/ba-version-2/internal/usecase"
	"github.com/arendi-project/ba-version-2/pkg/logger"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Feature struct {
	Carting usecase.Carting
	Order   usecase.Order
	User    usecase.User
}

type Middleware struct {
	Authentication middleware.Authentication
	Authorization  middleware.Authorization
}

func NewRouter(handler *gin.Engine, log logger.Interface, f *Feature, m *Middleware) {
	handler.Use(gin.Logger())
	handler.Use(gin.Recovery())
	// Swagger
	// TODO if you need this see https://github.com/evrone/go-clean-template

	// K8s probe
	handler.GET("/ping", func(c *gin.Context) { c.JSON(http.StatusOK, "pong") })
	handler.HEAD("/ping", func(c *gin.Context) { c.Status(http.StatusOK) })
	// Prometheus metrics
	// TODO if you need this see https://github.com/evrone/go-clean-template

	v1 := handler.Group("/v1", m.Authentication.Authenticate)
	{
		newCartingRoutes(v1, m.Authorization, f.Carting, log)
		newOrderRoutes(v1, m.Authorization, f.Order, log)
	}
}
