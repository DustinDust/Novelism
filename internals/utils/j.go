package utils

import (
	"errors"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
)

type JwtClaims struct {
	Claims interface{} `json:"claims"`
	jwt.RegisteredClaims
}

func SignToken(claims interface{}) (string, error) {
	secret := viper.GetViper().GetString("jwt.secret")
	expirationDuration := viper.GetDuration("jwt.duration")

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, JwtClaims{
		Claims: claims,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expirationDuration)),
		},
	})
	signed, err := token.SignedString([]byte(secret))
	return signed, err
}

func RetreiveUserIdFromContext(c echo.Context) (int, error) {
	log.Println(c.Get("user"))
	userId, ok := c.Get("user").(int)
	if !ok {
		return -1, errors.New("invalid server context")
	}

	return userId, nil
}
