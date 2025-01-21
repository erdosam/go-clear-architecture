package v1

import (
	"github.com/arendi-project/ba-version-2/internal/controller/http/middleware"
	"github.com/arendi-project/ba-version-2/internal/entity"
	"github.com/arendi-project/ba-version-2/internal/usecase"
	"github.com/arendi-project/ba-version-2/pkg/logger"
	"github.com/gin-gonic/gin"
	"net/http"
)

type cartingRoutes struct {
	usecase usecase.Carting
}

func newCartingRoutes(handler *gin.RouterGroup, m middleware.Authorization, cart usecase.Carting, log logger.Interface) {
	route := &cartingRoutes{
		usecase: cart,
	}

	h := handler.Group("/cart")
	{
		h.HEAD("/items", func(c *gin.Context) { c.Status(http.StatusOK) })
		h.GET("/items", m.Authorize("item", "read"), route.getCartItems)
		h.GET("/item/:id", m.Authorize("item", "read"), route.getCartItem)
		h.POST("/add-item", m.Authorize("item", "write"), route.addItemToCart)
	}
	log.Info("Done route : carting")
}

func (r *cartingRoutes) getCartItems(c *gin.Context) {
	cart, _ := r.usecase.GetCart(getIdentity(c))
	items, err := r.usecase.GetItems(cart)
	if err != nil {
		c.JSON(http.StatusNotFound, err)
		return
	}
	c.JSON(http.StatusOK, entity.CartItemResponse{Items: items})
}

func (r *cartingRoutes) getCartItem(c *gin.Context) {
	cart, _ := r.usecase.GetCart(getIdentity(c))
	id := c.Param("id")
	item, err := r.usecase.GetItem(cart, id)
	if err != nil {
		c.JSON(http.StatusNotFound, err.Error())
		return
	}
	c.JSON(http.StatusOK, item)
}

func (r *cartingRoutes) addItemToCart(c *gin.Context) {

}
