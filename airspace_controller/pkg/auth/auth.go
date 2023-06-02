package auth

import (
	"github.com/golang-jwt/jwt/v5"
	"time"
)

var JwtKey = []byte("SECRET_KEY")

type Credentials struct {
	Password string `json:"password"`
	Username string `json:"username"`
}

type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

type JWTGeneratedToken struct {
	Token      string    `json:"token"`
	Expiration time.Time `json:"expiration"`
}

func Authenticate(creds Credentials) (error, JWTGeneratedToken) {
	// TODO: Validate password

	expirationTime := time.Now().Add(1 * time.Hour)
	claims := &Claims{
		Username: creds.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(JwtKey)
	if err != nil {
		return err, JWTGeneratedToken{}
	}
	return nil, JWTGeneratedToken{
		Token:      tokenString,
		Expiration: expirationTime,
	}
}
