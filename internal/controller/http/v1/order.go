package v1

import (
	"github.com/arendi-project/ba-version-2/internal/controller/http/middleware"
	"github.com/arendi-project/ba-version-2/internal/usecase"
	"github.com/arendi-project/ba-version-2/pkg/logger"
	"github.com/gin-gonic/gin"
	"net/http"
)

type orderRoutes struct {
	usecase usecase.Order
}

func newOrderRoutes(handler *gin.RouterGroup, m middleware.Authorization, uc usecase.Order, log logger.Interface) {
	route := &orderRoutes{uc}
	h := handler.Group("/orders")
	{
		h.GET("active", m.Authorize("order", "list"), route.activeOrders)
	}
	g := handler.Group("/order")
	{
		g.POST("create", m.Authorize("order", "read"), route.orderView)
		g.GET("view", m.Authorize("order", "read"), route.orderView)
		g.POST("submit", m.Authorize("order", "read"), route.orderView)
	}
	log.Info("Done route : order")
}

func (r orderRoutes) activeOrders(c *gin.Context) {
	c.JSON(http.StatusOK, []interface{}{})
}

func (r orderRoutes) orderView(c *gin.Context) {
	order, err := r.usecase.GetOrderById("id")
	if err != nil {
		return
	}
	c.JSON(http.StatusOK, order)
}
