package middleware

import (
	"crypto/rsa"
	"github.com/arendi-project/ba-version-2/internal/entity"
	"github.com/arendi-project/ba-version-2/internal/usecase"
	"github.com/arendi-project/ba-version-2/pkg/logger"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
)

type authMiddleware struct {
	log     logger.Interface
	user    usecase.User
	keyFile string
}

func NewJwtAuthentication(user usecase.User, log logger.Interface, keyFile string) Authentication {
	return &authMiddleware{
		log:     log,
		user:    user,
		keyFile: keyFile,
	}
}

func (auth *authMiddleware) Authenticate(c *gin.Context) {
	tokenString := c.GetHeader("X-Ba-Token")
	token, err := jwt.ParseWithClaims(tokenString, &entity.CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			auth.log.Debug("Failed on token.Method.(*jwt.SigningMethodRSA)")
			return nil, http.ErrAbortHandler
		}
		return getRSAPublicKey(auth.keyFile)
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
		user, err := auth.user.GetUserById(claims.AccountId)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization"})
			c.Abort()
			return
		}
		c.Set("identity", user)
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
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
