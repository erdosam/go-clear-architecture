package middleware

import "github.com/gin-gonic/gin"

type (
	Authentication interface {
		Authenticate(c *gin.Context)
	}
	Authorization interface {
		Authorize(obj string, act string) gin.HandlerFunc
	}
)
