package utils

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
)

type JWTUtils struct{}

type JwtClaims struct {
	Claims interface{} `json:"claims"`
	jwt.RegisteredClaims
}

type JwtSignOption struct {
	ExpirationDuration time.Duration
}

type SignedJwtResult struct {
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expiresAt"`
}

func (j JWTUtils) SignToken(claims interface{}, option *JwtSignOption) (SignedJwtResult, error) {
	secret := viper.GetViper().GetString("jwt.secret")
	var expirationDuration time.Duration
	if option == nil || option.ExpirationDuration == 0 {
		expirationDuration = viper.GetViper().GetDuration("jwt.duration")
	} else {
		expirationDuration = option.ExpirationDuration
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, JwtClaims{
		Claims: claims,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expirationDuration)),
		},
	})
	signed, err := token.SignedString([]byte(secret))
	return SignedJwtResult{
		Token:     signed,
		ExpiresAt: time.Now().Add(expirationDuration),
	}, err
}

func (j JWTUtils) RetreiveUserIdFromContext(c echo.Context) (int, error) {
	userId, ok := c.Get("user").(int)
	if !ok {
		return -1, errors.New("invalid server context")
	}

	return userId, nil
}

func (j JWTUtils) VerifyToken(tokenString string) (jwt.MapClaims, error) {

	token, err := jwt.ParseWithClaims(tokenString, jwt.MapClaims{}, func(t *jwt.Token) (interface{}, error) {
		return viper.GetViper().GetString("jwt.secret"), nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return token.Claims.(jwt.MapClaims), nil
}

var JWT = JWTUtils{}
