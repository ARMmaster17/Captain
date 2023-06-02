package auth

import (
	"github.com/ARMmaster17/Captain/airspace_controller/pkg/auth"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"strings"
	"time"
)

func HandleLoginPost(c *gin.Context) {
	var creds auth.Credentials
	err := c.ShouldBindJSON(&creds)
	if err != nil {
		c.String(http.StatusBadRequest, "")
		return
	}

	var token auth.JWTGeneratedToken
	err, token = auth.Authenticate(creds)
	if err != nil {
		c.String(http.StatusInternalServerError, "")
		return
	}

	c.JSON(http.StatusOK, token)
}

func HandleRefreshPost(c *gin.Context) {
	reqToken := c.GetHeader("Authorization")
	splitToken := strings.Split(reqToken, "Bearer ")
	tknStr := splitToken[1]
	claims := &auth.Claims{}
	tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
		return auth.JwtKey, nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			c.String(http.StatusUnauthorized, "")
			return
		}
		c.String(http.StatusBadRequest, "")
		return
	}
	if !tkn.Valid {
		c.String(http.StatusUnauthorized, "")
		return
	}

	// TODO: The above should be middleware for each authed request.
	if time.Until(claims.ExpiresAt.Time) > 30*time.Second {
		c.String(http.StatusBadRequest, "")
		return
	}
	expirationTime := time.Now().Add(5 * time.Minute)
	claims.ExpiresAt = jwt.NewNumericDate(expirationTime)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(auth.JwtKey)
	if err != nil {
		c.String(http.StatusInternalServerError, "")
		return
	}

	c.SetCookie("token", tokenString, 5*60, "/", "localhost", false, false)
	c.String(http.StatusOK, tokenString)
}
