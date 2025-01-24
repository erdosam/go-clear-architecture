package middleware

import (
	"github.com/arendi-project/ba-version-2/pkg/logger"
	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"net/http"
)

type casbinAuthorization struct {
	log      logger.Interface
	enforcer *casbin.Enforcer
}

func NewAbacAuthorization(l logger.Interface, e *casbin.Enforcer) Authorization {
	return &casbinAuthorization{l, e}
}

func (acc *casbinAuthorization) Authorize(obj string, act string) gin.HandlerFunc {
	return func(c *gin.Context) {
		acc.log.Debug("Authorizing %s %s", act, obj)
		if ok, _ := acc.enforcer.Enforce(c.MustGet("identity"), "obj", "act"); !ok {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}
		c.Next()
	}
}
