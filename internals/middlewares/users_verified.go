package middlewares

import (
	"fmt"
	"gin_stuff/internals/models"
	"gin_stuff/internals/utils"

	"github.com/labstack/echo/v4"
)

// need the model object
func NewUserVerificationRequireMiddleware(model models.UserRepository) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			userId, err := utils.JWT.RetreiveUserIdFromContext(c)
			if err != nil {
				utils.Logger.Err(fmt.Errorf("can't find user id inside context object: %v", err)).Send()
				return utils.ErrorUnauthorized
			}
			user, err := model.Get(int64(userId))
			if err != nil {
				utils.Logger.Err(fmt.Errorf("can't find user: %v", err))
				return utils.ErrorUnauthorized
			}
			if !user.Verified {
				return echo.NewHTTPError(utils.ErrorForbiddenResource.Code, utils.ErrorForbiddenResource)
			}
			return next(c)
		}
	}
}
