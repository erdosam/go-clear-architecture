package middleware

import (
	"github.com/arendi-project/ba-version-2/pkg/logger"
	"github.com/gin-gonic/gin"
)

type casbinAuthorization struct {
	log logger.Interface
}

func NewAbacAuthorization(l logger.Interface) Authorization {
	return &casbinAuthorization{l}
}

func (casbin *casbinAuthorization) Authorize(obj string, act string) gin.HandlerFunc {
	return func(c *gin.Context) {
		casbin.log.Debug("Authorizing %s %s", act, obj)
		c.Next()
	}
}
