package middleware

import (
	"github.com/ARMmaster17/Captain/airspace_controller/pkg/auth"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"regexp"
	"strings"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		bearerToken := c.GetHeader("Authorization")
		c.FullPath()
		isApiRoute, _ := regexp.MatchString("/api/.*", c.Request.URL.RequestURI())
		if !isApiRoute {
			c.Next()
			return
		}
		isAuthRoute, _ := regexp.MatchString("/api/v[0-9]+/auth/[a-z]+", c.Request.URL.RequestURI())
		if isAuthRoute {
			c.Next()
			return
		}
		if bearerToken == "" {
			c.String(http.StatusForbidden, "Missing authentication token")
			return
		}
		token := strings.ReplaceAll(bearerToken, "Bearer ", "")
		claims := &auth.Claims{}
		tkn, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
			return auth.JwtKey, nil
		})
		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				c.String(http.StatusUnauthorized, "Invalid signature")
				return
			}
			c.String(http.StatusBadRequest, "Malformed signature")
			return
		}
		if !tkn.Valid {
			c.String(http.StatusUnauthorized, "Invalid token")
			return
		}
		c.Next()
	}
}
