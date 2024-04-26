package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
)

type JWTService struct{}

type JwtClaims struct {
	Claims interface{} `json:"claims"`
	jwt.RegisteredClaims
}

type JwtSignOption struct {
	Secret             string
	ExpirationDuration time.Duration
	SigningMethod      jwt.SigningMethod
}

type SignedJwtResult struct {
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expiresAt"`
}

func (j JWTService) signToken(claims interface{}, option *JwtSignOption) (SignedJwtResult, error) {
	expiresAt := time.Now().Add(option.ExpirationDuration)
	token := jwt.NewWithClaims(option.SigningMethod, JwtClaims{
		Claims: claims,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
		},
	})
	signed, err := token.SignedString([]byte(option.Secret))
	return SignedJwtResult{
		Token:     signed,
		ExpiresAt: expiresAt,
	}, err
}

func (j JWTService) SignAccessToken(claims interface{}) (SignedJwtResult, error) {
	secret := viper.GetViper().GetString("jwt.secret")
	expirationDuration := viper.GetViper().GetDuration("jwt.duration")

	return j.signToken(claims, &JwtSignOption{
		Secret:             secret,
		ExpirationDuration: expirationDuration,
		SigningMethod:      jwt.SigningMethodHS256,
	})
}

func (j JWTService) SignRefreshToken(claims interface{}) (SignedJwtResult, error) {
	secret := viper.GetViper().GetString("jwt.refreshSecret")
	expirationDuration := viper.GetViper().GetDuration("jwt.refreshDuration")

	return j.signToken(claims, &JwtSignOption{
		Secret:             secret,
		ExpirationDuration: expirationDuration,
		SigningMethod:      jwt.SigningMethodHS256,
	})
}

func (j JWTService) RetreiveUserIdFromContext(c echo.Context) (int, error) {
	userId, ok := c.Get("user").(int)
	if !ok {
		return -1, errors.New("invalid server context")
	}

	return userId, nil
}

func (j JWTService) VerifyAccessToken(tokenString string) (jwt.MapClaims, error) {

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
