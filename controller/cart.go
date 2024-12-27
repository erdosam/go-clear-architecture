package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func CartList(c *gin.Context) {
	c.JSON(http.StatusOK, "")
}
