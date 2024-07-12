package services

import (
	"gin_stuff/internals/utils"
	"time"

	"github.com/golang-jwt/jwt/v5"
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

func (j JWTService) verifyToken(tokenString string, secret string) (jwt.MapClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, jwt.MapClaims{}, func(t *jwt.Token) (interface{}, error) {
		return secret, nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, utils.ErrorInvalidToken()
	}
	return token.Claims.(jwt.MapClaims), nil
}

func (j JWTService) VerifyAccessToken(tokenString string) (jwt.MapClaims, error) {
	secret := viper.GetViper().GetString("jwt.accessSecret")
	return j.verifyToken(tokenString, secret)
}
