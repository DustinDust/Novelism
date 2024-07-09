package services

import (
	"gin_stuff/internals/utils"
	"net/http"
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
	secret := viper.GetViper().GetString("jwt.access_secret")
	expirationDuration := viper.GetViper().GetDuration("jwt.access_expiration_duration")

	return j.signToken(claims, &JwtSignOption{
		Secret:             secret,
		ExpirationDuration: expirationDuration,
		SigningMethod:      jwt.SigningMethodHS256,
	})
}

// unused
func (j JWTService) SignRefreshToken(claims interface{}) (SignedJwtResult, error) {
	secret := viper.GetViper().GetString("jwt.refresh_secret")
	expirationDuration := viper.GetViper().GetDuration("jwt.refresh_expiration_duration")

	return j.signToken(claims, &JwtSignOption{
		Secret:             secret,
		ExpirationDuration: expirationDuration,
		SigningMethod:      jwt.SigningMethodHS256,
	})
}

func (j JWTService) RetrieveUserIdFromContext(c echo.Context) (int, error) {
	userId, ok := c.Get("user").(int)
	if !ok {
		return -1, echo.NewHTTPError(http.StatusInternalServerError, utils.ErrInvalidContext())
	}

	return userId, nil
}

func (j JWTService) verifyToken(tokenString string, secret string) (jwt.MapClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, jwt.MapClaims{}, func(t *jwt.Token) (interface{}, error) {
		return secret, nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, echo.NewHTTPError(http.StatusUnauthorized, utils.ErrorInvalidToken())
	}
	return token.Claims.(jwt.MapClaims), nil
}

func (j JWTService) VerifyAccessToken(tokenString string) (jwt.MapClaims, error) {
	secret := viper.GetViper().GetString("jwt.accessSecret")
	return j.verifyToken(tokenString, secret)
}

// unused
func (j JWTService) VerifyRefreshToken(tokenString string) (jwt.MapClaims, error) {
	secret := viper.GetViper().GetString("jwt.refreshSecret")
	return j.verifyToken(tokenString, secret)
}
