package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Ping(c *gin.Context) {
	c.JSON(http.StatusOK, "pong")
}

func PingHead(c *gin.Context) {
	c.Status(http.StatusOK)
}
