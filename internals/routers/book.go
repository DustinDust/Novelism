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
		return r.badRequestError(err)
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
			return r.notFoundError(err)
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
		return r.badRequestError(err)
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
	if updateBookPayload.Title != "" {
		book.Title = updateBookPayload.Title
	}
	if updateBookPayload.Description != "" {
		book.Description = updateBookPayload.Description
	}
	err = r.Model.Book.Update(book)
	if err != nil {
		r.serverError(err)
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
			return r.notFoundError(err)
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

func (r Router) FindBooks(c echo.Context) error {
	currentUserId, err := utils.RetreiveUserIdFromContext(c)
	if err != nil {
		return r.unauthorizedError(err)
	}
	filter := models.Filter{
		Page:         1,
		PageSize:     10,
		SortSafeList: []string{"id", "title", "-id", "-title", "created_at", "-created_at"},
		Sort:         "created_at", // default sort
	}

	userId := currentUserId // default to current session
	var title string

	queryParams := c.QueryParams()
	if queryParams.Has("userId") {
		var err error
		userId, err = strconv.Atoi(queryParams.Get("userId"))
		if err != nil {
			return r.badRequestError(err)
		}
	}
	if queryParams.Has("title") {
		title = queryParams.Get("title")
	}
	if queryParams.Has("page") {
		page, err := strconv.Atoi(queryParams.Get("page"))
		if err != nil {
			return r.badRequestError(err)
		}
		filter.Page = page
	}
	if queryParams.Has("pageSize") {
		pageSize, err := strconv.Atoi(queryParams.Get("pageSize"))
		if err != nil {
			return r.badRequestError(err)
		}
		filter.PageSize = pageSize
	}
	if queryParams.Has("sort") {
		filter.Sort = queryParams.Get("sort")
	}
	books, metadata, err := r.Model.Book.Find(userId, title, filter)
	if err != nil {
		return r.serverError(err)
	}
	return c.JSON(http.StatusOK,
		echo.Map{
			"metadata": metadata,
			"data": echo.Map{
				"books": books,
			},
			"ok": true,
		},
	)
}
