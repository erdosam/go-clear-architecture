package middleware

import (
	pgadapter "github.com/casbin/casbin-pg-adapter"
	"github.com/casbin/casbin/v2"
	"github.com/erdosam/go-clear-architecture/config"
	"github.com/erdosam/go-clear-architecture/internal/entity"
	"github.com/erdosam/go-clear-architecture/pkg/logger"
	"github.com/gin-gonic/gin"
	"net/http"
)

type casbinAuthorization struct {
	log      logger.Interface
	enforcer *casbin.Enforcer
}

func NewAbacAuthorization(cfg *config.Config, l logger.Interface) Authorization {
	pga, err := pgadapter.NewAdapter(cfg.PG.URL)
	if err != nil {
		l.Fatal(err)
	}
	enforcer, err := casbin.NewEnforcer(cfg.Casbin.ModelFile, pga)
	if err != nil {
		l.Fatal(err)
	}
	err = enforcer.LoadPolicy()
	if err != nil {
		l.Fatal(err)
	}
	return &casbinAuthorization{l, enforcer}
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
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Unauthorized action"})
			return
		}
		c.Next()
	}
}

func (acc *casbinAuthorization) LoadPolicy(c *gin.Context) {
	acc.log.Info("Load policy")
	err := acc.enforcer.LoadPolicy()
	if err != nil {
		acc.log.Error("Error load policy: %s", err.Error())
	}
	c.Status(http.StatusOK)
}
