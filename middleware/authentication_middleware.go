package middleware

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/snck/book-keeper-api/service"
)

type AuthenticationMiddleware struct {
	service *service.AuthenticationService
}

func NewAuthenticationMiddleware(service *service.AuthenticationService) *AuthenticationMiddleware {
	return &AuthenticationMiddleware{service: service}
}

func (m *AuthenticationMiddleware) ValidateToken() gin.HandlerFunc {
	return func(c *gin.Context) {

		token := c.GetHeader("Authorization")
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "no authorization header"})
			c.Abort()
			return
		}

		token = strings.TrimPrefix(token, "Bearer ")
		claims, err := m.service.ParseAndValidateToken(token)

		if err != nil {
			if errors.Is(err, jwt.ErrTokenExpired) || errors.Is(err, jwt.ErrTokenSignatureInvalid) {
				c.JSON(http.StatusUnauthorized, gin.H{"message": "invalid token"})
				c.Abort()
				return
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"message": "error parsing token"})
				c.Abort()
				return
			}
		}

		if claims == nil {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "invalid token"})
			c.Abort()
			return
		}

		blocked, err := m.service.IsTokenExistInBlocklist(token)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "error with database"})
			c.Abort()
			return
		}

		if blocked {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "invalid token"})
			c.Abort()
			return
		}

		c.Set("userID", claims.ID)

		c.Next()
	}
}
