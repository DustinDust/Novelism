package middlewares

import (
	"errors"
	"fmt"

	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
)

// @param{kind}: "access" | "refresh"
func NewJWTMiddleware(kind string) echo.MiddlewareFunc {
	jwtSecret := viper.GetViper().GetString(fmt.Sprintf("jwt.%sSecret", kind))

	return echojwt.WithConfig(echojwt.Config{
		ContextKey:  "user",
		TokenLookup: "header:Authorization:Bearer ",
		SigningKey:  []byte(jwtSecret),
		ParseTokenFunc: func(c echo.Context, auth string) (interface{}, error) {
			customClaims := new(struct {
				Claims int `json:"claims"`
				jwt.RegisteredClaims
			})
			token, err := jwt.ParseWithClaims(auth, customClaims, func(t *jwt.Token) (interface{}, error) {
				if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, errors.New("unexpected signing method")
				}
				return []byte(jwtSecret), nil
			})
			if err != nil {
				return nil, err
			}
			content, ok := token.Claims.(*struct {
				Claims int `json:"claims"`
				jwt.RegisteredClaims
			})
			if !ok {
				return nil, errors.New("unexpected jwt format")
			}
			return content.Claims, nil
		},
	})
}
