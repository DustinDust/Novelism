package router

import (
	"gin_stuff/internals/models"
	"gin_stuff/internals/utils"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
)

type Router struct {
	Model *models.Models
}

func NewRouter(model *models.Models) Router {
	return Router{
		Model: model,
	}
}

func (r Router) GetConfig(c echo.Context) error {
	key := c.Param("key")
	userId, err := utils.RetreiveUserIdFromContext(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	user, err := r.Model.User.Get(int64(userId))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	if len(key) <= 0 {
		return echo.NewHTTPError(http.StatusBadRequest)
	}
	return c.JSON(http.StatusOK, echo.Map{
		"user": user,
		key:    viper.GetViper().Get(key),
	})
}
