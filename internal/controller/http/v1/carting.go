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
	r := &cartingRoutes{
		usecase: rh.feature.Carting,
	}
	h := parent.Group("/cart")
	{
		h.HEAD("/items", func(c *gin.Context) { c.Status(http.StatusOK) })
		h.GET("/items", rh.access.Authorize("cart-item", "read"), r.getCartItems)
		h.GET("/item/:id", rh.access.Authorize("cart-item", "read"), r.getCartItem)
		h.POST("/item/add", rh.access.Authorize("cart-item", "add"), r.addItemToCart)
		h.POST("/item/edit/:id", r.editCartItemResource, rh.access.Authorize("cart-item", "edit"), r.editCartItem)
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
	var body, err = shouldBindJSON[entity.AddItemToCartForm](c)
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

func (r *cartingRoutes) editCartItemResource(c *gin.Context) {
	cart, _ := r.usecase.GetCart(getIdentity(c))
	item, _ := r.usecase.GetItem(cart, c.Param("id"))
	setAccessResource(c, item)
}

func (r *cartingRoutes) editCartItem(c *gin.Context) {
	var body, err = shouldBindJSON[entity.EditCartItemForm](c)
	if err != nil {
		return
	}
	err = r.usecase.EditCartItem(body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, "")
}
