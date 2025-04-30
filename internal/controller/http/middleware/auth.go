package middleware

import (
	"crypto/rsa"
	"github.com/dgrijalva/jwt-go"
	"github.com/erdosam/go-clear-architecture/config"
	"github.com/erdosam/go-clear-architecture/internal/entity"
	"github.com/erdosam/go-clear-architecture/internal/usecase"
	"github.com/erdosam/go-clear-architecture/pkg/logger"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
)

type authMiddleware struct {
	log    logger.Interface
	user   usecase.User
	config *config.Config
}

func NewJwtAuthentication(user usecase.User, log logger.Interface, cfg *config.Config) Authentication {
	return &authMiddleware{
		log:    log,
		user:   user,
		config: cfg,
	}
}

func (auth *authMiddleware) Authenticate(c *gin.Context) {
	tokenString := c.GetHeader(auth.config.HTTP.AuthKey)
	token, err := jwt.ParseWithClaims(tokenString, &entity.CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			auth.log.Debug("Failed on token.Method.(*jwt.SigningMethodRSA)")
			return nil, http.ErrAbortHandler
		}
		return getRSAPublicKey(auth.config.ClientKeyFile)
	})

	if err != nil || !token.Valid {
		auth.log.Error(err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		c.Abort()
		return
	}
	// Set the token claims to the context
	if claims, ok := token.Claims.(*entity.CustomClaims); ok && token.Valid {
		auth.log.Debug("Account Id: %s", claims.AccountId)
		user, err := auth.user.GetUserFromId(claims.AccountId, claims.ClientKey)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authentication"})
			c.Abort()
			return
		}
		c.Set(IdentityContextKey, user)
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthenticated"})
		c.Abort()
		return
	}
	c.Next()
}

func getRSAPublicKey(filePath string) (*rsa.PublicKey, error) {
	var key *rsa.PublicKey
	verifyBytes, err := os.ReadFile(filePath)
	if err != nil {
		return key, err
	}
	return jwt.ParseRSAPublicKeyFromPEM(verifyBytes)
}
