package middleware

import (
	"github.com/arendi-project/ba-version-2/internal/entity"
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

func (acc *casbinAuthorization) Authorize(dom string, act string) gin.HandlerFunc {
	return func(c *gin.Context) {
		acc.log.Debug("Authorizing %s %s", act, dom)
		user := c.MustGet(IdentityContextKey).(entity.User)
		obj, exists := c.Get(ResourceContextKey)
		if !exists {
			obj = ""
		}
		if ok, _ := acc.enforcer.Enforce(user, obj, dom, act); !ok {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"message": "Unauthorized action"})
			return
		}
		c.Next()
	}
}
