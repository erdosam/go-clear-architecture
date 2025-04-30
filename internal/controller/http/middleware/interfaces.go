package middleware

import "github.com/gin-gonic/gin"

const (
	IdentityContextKey = "identity"
	ResourceContextKey = "casbinObj"
)

type (
	Authentication interface {
		Authenticate(c *gin.Context)
	}
	Authorization interface {
		LoadPolicy(c *gin.Context)
		Authorize(dom string, act string) gin.HandlerFunc
	}
)
