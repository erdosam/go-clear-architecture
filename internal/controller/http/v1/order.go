package v1

import (
	"github.com/arendi-project/ba-version-2/internal/usecase"
	"github.com/gin-gonic/gin"
	"net/http"
)

type orderRoutes struct {
	usecase usecase.Order
}

func (rh *routerHandler) initOrderRoutes(parent *gin.RouterGroup) {
	route := &orderRoutes{rh.feature.Order}
	m := rh.access
	h := parent.Group("/orders")
	{
		h.GET("active", m.Authorize("order", "list"), route.activeOrders)
	}
	g := parent.Group("/order")
	{
		g.POST("create", m.Authorize("order", "read"), route.orderView)
		g.GET("view", m.Authorize("order", "read"), route.orderView)
		g.POST("submit", m.Authorize("order", "read"), route.orderView)
	}
	rh.log.Info("Done route : order")
}

func (r *orderRoutes) activeOrders(c *gin.Context) {
	c.JSON(http.StatusOK, []interface{}{})
}

func (r *orderRoutes) orderView(c *gin.Context) {
	order, err := r.usecase.GetOrderById("id")
	if err != nil {
		return
	}
	c.JSON(http.StatusOK, order)
}
