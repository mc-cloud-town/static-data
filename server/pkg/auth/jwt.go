package auth

import (
	"errors"
	"time"

	"server/config"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

const JWT_ISSUER = "ctec-api-server"

var ErrTokenValidation = errors.New("token validation failed")

type JwtCustomClaims struct {
	jwt.RegisteredClaims
	ID uint `json:"id"`
}

type JwtToken struct {
	Token  string          `json:"token"`
	Claims JwtCustomClaims `json:"claims"`
}

func New() *JwtToken {
	return &JwtToken{}
}

// CreateUserToken creates a token for the user
func CreateUserToken(ID uint) (*JwtToken, error) {
	expiresAt := time.Now().Add(time.Hour * 24 * time.Duration(config.Get().API.JwtExpireDay))
	claims := JwtCustomClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    JWT_ISSUER,
			ExpiresAt: jwt.NewNumericDate(expiresAt),
		},
		ID: ID,
	}

	t := jwt.New(jwt.SigningMethodES256)
	t.Claims = &claims
	tokenString, err := t.SignedString(config.GetJwtPrivateKey())
	if err != nil {
		return nil, err
	}

	return &JwtToken{Token: tokenString, Claims: claims}, nil
}

// ParseToken parses the token and returns the claims
func (j *JwtToken) ParseToken(token string) (*JwtCustomClaims, error) {
	t, err := jwt.ParseWithClaims(token, &JwtCustomClaims{}, func(t *jwt.Token) (any, error) {
		return config.GetJwtPublicKey(), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := t.Claims.(*JwtCustomClaims)
	if !ok || !t.Valid {
		return nil, ErrTokenValidation
	}

	return claims, nil
}

// JwtClaims returns the claims from the token
func (j *JwtToken) JwtClaims(c *gin.Context) (*JwtCustomClaims, error) {
	token := c.GetHeader("Authorization")
	claims, err := j.ParseToken(token)
	return claims, err
}

// JwtUserId returns the user ID from the token
func (j *JwtToken) JwtUserId(c *gin.Context) uint {
	claims, _ := j.JwtClaims(c)
	return claims.ID
}
