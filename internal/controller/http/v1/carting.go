package v1

import (
	"github.com/arendi-project/ba-version-2/internal/entity"
	"github.com/arendi-project/ba-version-2/internal/usecase"
	"github.com/gin-gonic/gin"
	"net/http"
)

type cartingRoutes struct {
	usecase usecase.Carting
}

func (rh *routerHandler) initCartingRoutes(parent *gin.RouterGroup) {
	route := &cartingRoutes{
		usecase: rh.feature.Carting,
	}
	h := parent.Group("/cart")
	{
		h.HEAD("/items", func(c *gin.Context) { c.Status(http.StatusOK) })
		h.GET("/items", rh.access.Authorize("item", "read"), route.getCartItems)
		h.GET("/item/:id", rh.access.Authorize("item", "read"), route.getCartItem)
		h.POST("/add-item", rh.access.Authorize("item", "write"), route.addItemToCart)
	}
	rh.log.Info("Done route : carting")
}

func (r *cartingRoutes) getCartItems(c *gin.Context) {
	cart, _ := r.usecase.GetCart(getIdentity(c))
	items, err := r.usecase.GetItems(cart)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, entity.CartItemResponse{Items: items})
}

func (r *cartingRoutes) getCartItem(c *gin.Context) {
	cart, _ := r.usecase.GetCart(getIdentity(c))
	id := c.Param("id")
	item, err := r.usecase.GetItem(cart, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, item)
}

func (r *cartingRoutes) addItemToCart(c *gin.Context) {
	var body, err = shouldBindJSON[entity.AddItemToCartRequest](c)
	if err != nil {
		return
	}
	err = r.usecase.AddItemToCart(body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, "")
}
