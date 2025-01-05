package middleware

import (
	"github.com/arendi-project/ba-version-2/internal/entity"
	"github.com/arendi-project/ba-version-2/internal/usecase"
	"github.com/arendi-project/ba-version-2/pkg/logger"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
)

type authMiddleware struct {
	log     logger.Interface
	user    usecase.User
	keyFile string
}

var auth *authMiddleware

func NewJwtAuthentication(handler *gin.RouterGroup, user usecase.User, log logger.Interface, keyFile string) {
	auth = &authMiddleware{
		log:     log,
		user:    user,
		keyFile: keyFile,
	}
	handler.Use(jwtAuth)
}

func jwtAuth(c *gin.Context) {
	tokenString := c.GetHeader("X-Ba-Token")
	token, err := jwt.ParseWithClaims(tokenString, &entity.CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, http.ErrAbortHandler
		}
		return []byte(auth.keyFile), nil
	})
	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		c.Abort()
		return
	}
	// Set the token claims to the context
	if claims, ok := token.Claims.(entity.CustomClaims); ok && token.Valid {
		auth.log.Debug(claims)
		_, err := auth.user.GetUserById(claims.AccountId)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization"})
			return
		}
		c.Set("identity", claims)
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		c.Abort()
		return
	}
	c.Next()
}
