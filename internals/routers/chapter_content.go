package router

import (
	"errors"
	"gin_stuff/internals/repositories"
	"gin_stuff/internals/utils"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

func (r Router) GetContent(e echo.Context) error {
	chapterIdStr := e.Param("chapterId")
	chapterId, err := strconv.Atoi(chapterIdStr)
	if err != nil {
		return r.badRequestError(utils.ErrorInvalidRouteParam)
	}

	content, err := r.Repository.Content.Get(int64(chapterId))
	if err != nil {
		switch {
		case errors.Is(err, utils.ErrorRecordsNotFound):
			return r.notFoundError(err)
		default:
			return r.serverError(err)
		}
	}

	return e.JSON(http.StatusOK, Response[repositories.Content]{
		OK:   true,
		Data: *content,
	})
}
