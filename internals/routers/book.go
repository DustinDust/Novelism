package router

import (
	"errors"
	"gin_stuff/internals/models"
	"gin_stuff/internals/utils"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type CreateBookPayload struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type UpdateBookPayload struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

func (r Router) CreateBook(c echo.Context) error {
	validate := utils.NewValidator()
	createBookPayload := new(CreateBookPayload)
	if err := c.Bind(createBookPayload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if err := validate.ValidateStruct(createBookPayload); err != nil {
		if verr, ok := err.(*utils.StructValidationErrors); ok {
			return verr.TranslateError()
		} else {
			return r.serverError(err)
		}
	}
	userId, err := utils.RetreiveUserIdFromContext(c)
	if err != nil {
		return r.forbiddenError(err)
	}
	user, err := r.Model.User.Get(int64(userId))
	if err != nil {
		return r.serverError(err)
	}
	book := models.Book{
		UserID:      user.ID,
		User:        user,
		Title:       createBookPayload.Title,
		Description: createBookPayload.Description,
	}
	if err := r.Model.Book.Insert(&book); err != nil {
		return r.badRequestError(err)
	}
	return c.JSON(http.StatusCreated, echo.Map{
		"ok":   true,
		"data": book,
	})
}

func (r Router) GetBook(c echo.Context) error {
	idStr := c.Param("id")
	if idStr == "" {
		return r.badRequestError(utils.ErrorInvalidRouteParam)
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return r.badRequestError(err)
	}
	userId, err := utils.RetreiveUserIdFromContext(c)
	if err != nil {
		return r.forbiddenError(err)
	}
	book, err := r.Model.Book.Get(int64(id))
	if err != nil {
		switch {
		case errors.Is(err, utils.ErrorRecordsNotFound):
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		default:
			return r.serverError(err)
		}
	}
	if book.UserID != int64(userId) {
		return r.forbiddenError(utils.ErrorForbiddenResource)
	}
	return c.JSON(http.StatusOK, echo.Map{
		"ok":   true,
		"data": book,
	})
}

func (r Router) UpdateBook(c echo.Context) error {
	idStr := c.Param("id")
	if idStr == "" {
		return r.badRequestError(utils.ErrorInvalidRouteParam)
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return r.badRequestError(err)
	}

	validate := utils.NewValidator()
	updateBookPayload := new(UpdateBookPayload)
	if err := c.Bind(updateBookPayload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if err := validate.ValidateStruct(updateBookPayload); err != nil {
		if verr, ok := err.(*utils.StructValidationErrors); ok {
			return verr.TranslateError()
		} else {
			return r.serverError(err)
		}
	}
	userId, err := utils.RetreiveUserIdFromContext(c)
	if err != nil {
		return r.forbiddenError(err)
	}
	book, err := r.Model.Book.Get(int64(id))
	if err != nil {
		switch {
		case errors.Is(err, utils.ErrorRecordsNotFound):
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		default:
			return r.serverError(err)
		}
	}
	if book.UserID != int64(userId) {
		return r.forbiddenError(utils.ErrorForbiddenResource)
	}
	if book.Title != "" {
		book.Title = updateBookPayload.Title
	}
	if book.Description != "" {
		book.Description = updateBookPayload.Description
	}
	err = r.Model.Book.Update(book)
	if err != nil {
		r.badRequestError(err)
	}
	return c.JSON(http.StatusOK, echo.Map{"ok": true, "data": book})
}

func (r Router) DeleteBook(c echo.Context) error {
	idStr := c.Param("id")
	if idStr == "" {
		return r.badRequestError(utils.ErrorInvalidRouteParam)
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return r.badRequestError(err)
	}
	userId, err := utils.RetreiveUserIdFromContext(c)
	if err != nil {
		return r.forbiddenError(err)
	}
	book, err := r.Model.Book.Get(int64(id))
	if err != nil {
		switch {
		case errors.Is(err, utils.ErrorRecordsNotFound):
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		default:
			return r.serverError(err)
		}
	}
	if book.UserID != int64(userId) {
		return r.forbiddenError(utils.ErrorForbiddenResource)
	}
	err = r.Model.Book.Delete(int64(id))
	if err != nil {
		return r.serverError(err)
	}
	return c.JSON(http.StatusOK, echo.Map{"ok": true})
}
