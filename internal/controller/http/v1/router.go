package v1

import (
	"github.com/arendi-project/ba-version-2/internal/controller/http/middleware"
	"github.com/arendi-project/ba-version-2/internal/entity"
	"github.com/arendi-project/ba-version-2/internal/usecase"
	"github.com/arendi-project/ba-version-2/pkg/logger"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Feature struct {
	Carting  usecase.Carting
	Order    usecase.Order
	User     usecase.User
	Category usecase.Category
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
	// Prometheus metrics
	// TODO if you need this see https://github.com/evrone/go-clean-template

	v1 := handler.Group("/v1", m.Authentication.Authenticate)
	{
		handler.initCategoryRoutes(v1)
		handler.initCartingRoutes(v1)
		handler.initOrderRoutes(v1)
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

func shouldBindJSON[T any](c *gin.Context) (T, error) {
	var body T
	err := c.ShouldBindJSON(&body)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusPreconditionFailed, gin.H{"error": "Invalid in type a json field"})
	}
	return body, err
}
