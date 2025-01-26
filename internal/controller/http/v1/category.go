package v1

import (
	"github.com/arendi-project/ba-version-2/internal/usecase"
	"github.com/gin-gonic/gin"
	"net/http"
)

type categoryRoute struct {
	usecase usecase.Category
}

func (rh *routerHandler) initCategoryRoutes(parent *gin.RouterGroup) {
	r := &categoryRoute{rh.feature.Category}

	parent.GET("/categories", r.getAvailableCategories)
	h := parent.Group("/category")
	{
		h.HEAD("/:id", func(c *gin.Context) { c.Status(http.StatusOK) })
	}
	rh.log.Info("Done route : category")
}

func (r *categoryRoute) getAvailableCategories(c *gin.Context) {
	categories, err := r.usecase.GetAvailableCategories()
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, categories)
}
