package v1

import (
	"github.com/arendi-project/ba-version-2/internal/entity"
	"github.com/arendi-project/ba-version-2/internal/usecase"
	"github.com/arendi-project/ba-version-2/pkg/logger"
	"github.com/gin-gonic/gin"
	"net/http"
)

type cartingRoutes struct {
	usecase usecase.Carting
}

func newCartingRoutes(handler *gin.RouterGroup, cart usecase.Carting, log logger.Interface) {
	uc := &cartingRoutes{
		usecase: cart,
	}

	h := handler.Group("/cart")
	{
		h.HEAD("/items", func(c *gin.Context) { c.Status(http.StatusOK) })
		h.GET("/items", uc.getItems)
		h.POST("/add-item", uc.addItem)
	}
	log.Info("Done route : carting")
}

type cartItemResponse struct {
	Items []entity.CartItem `json:"items"`
}

func (r *cartingRoutes) getItems(c *gin.Context) {
	cart := entity.Cart{UserId: c.Query("user_id")}
	items, err := r.usecase.GetItems(cart)
	if err != nil {
		c.JSON(http.StatusNotFound, err)
		return
	}
	c.JSON(http.StatusOK, cartItemResponse{items})
}

func (r *cartingRoutes) addItem(c *gin.Context) {

}
