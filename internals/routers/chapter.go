package router

import (
	"errors"
	"gin_stuff/internals/repositories"
	"gin_stuff/internals/utils"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type CreateChapterPayload struct {
	BookID      int    `validate:"required,gte=1"`
	Title       string `json:"title" validate:"required,max=128"`
	Description string `json:"description"`
}

type UpdateChapterPayload struct {
	Title       string `json:"title" validate:"max=128"`
	Description string `json:"description"`
}

func (r Router) CreateChapter(c echo.Context) error {
	validate := utils.NewValidator()
	createChapterPayload := new(CreateChapterPayload)
	idStr := c.Param("bookId")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return r.badRequestError(err)
	}
	userId, err := r.JwtService.RetrieveUserIdFromContext(c)
	if err != nil {
		return r.forbiddenError(err)
	}
	book, err := r.Repository.Book.Get(int64(id))
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
	createChapterPayload.BookID = id
	if err := c.Bind(createChapterPayload); err != nil {
		return r.badRequestError(err)
	}
	if err := validate.ValidateStruct(createChapterPayload); err != nil {
		if verr, ok := err.(*utils.StructValidationErrors); ok {
			return verr.TranslateError()
		} else {
			return r.serverError(err)
		}
	}
	chapter := repositories.Chapter{
		AuthorID:    book.UserID,
		BookID:      book.ID,
		Title:       createChapterPayload.Title,
		Description: createChapterPayload.Description,
	}
	err = r.Repository.Chapter.Insert(&chapter)
	if err != nil {
		return r.serverError(err)
	}
	return c.JSON(http.StatusCreated, Response[repositories.Chapter]{
		OK:   true,
		Data: chapter,
	})
}

// get all chapters from 1 book (with optional filter)
func (r Router) FindChapters(c echo.Context) error {
	bookIdStr := c.Param("bookId")
	bookId, err := strconv.Atoi(bookIdStr)
	if err != nil {
		return r.badRequestError(err)
	}
	_, err = r.Repository.Book.Get(int64(bookId))
	if err != nil {
		switch {
		case errors.Is(err, utils.ErrorRecordsNotFound):
			return r.notFoundError(err)
		default:
			return r.serverError(err)
		}
	}
	var title string
	filter := repositories.Filter{
		Page:     1,
		PageSize: 10,
		SortSafeList: []string{
			"id",
			"title",
			"-id",
			"-title",
			"chapter_no",
			"-chapter_no",
			"created_at",
			"-created_at",
			"updated_at",
			"-updated_at",
		},
		Sort: "chapter_no",
	}
	queryParams := c.QueryParams()
	if queryParams.Has("title") {
		title = queryParams.Get("title")
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
	if queryParams.Has("page") {
		page, err := strconv.Atoi(queryParams.Get("page"))
		if err != nil {
			return r.badRequestError(err)
		}
		filter.Page = page
	}

	chapters, metadata, err := r.Repository.Chapter.Find(int64(bookId), title, filter)
	if err != nil {
		return r.serverError(err)
	}
	return c.JSON(http.StatusOK, Response[[]*repositories.Chapter]{
		OK:       true,
		Metadata: metadata,
		Data:     chapters,
	})
}

func (r Router) UpdateChapter(c echo.Context) error {
	chapterNo, err := strconv.Atoi(c.Param("chapterNo"))
	if err != nil {
		return r.badRequestError(err)
	}
	bookId, err := strconv.Atoi(c.Param("bookId"))
	if err != nil {
		return r.badRequestError(err)
	}
	userId, err := r.JwtService.RetrieveUserIdFromContext(c)
	if err != nil {
		return r.unauthorizedError(err)
	}
	validate := utils.NewValidator()
	updateChapterPayload := new(UpdateChapterPayload)
	if err := c.Bind(updateChapterPayload); err != nil {
		return r.badRequestError(err)
	}
	if err := validate.ValidateStruct(updateChapterPayload); err != nil {
		if verr, ok := err.(*utils.StructValidationErrors); ok {
			return verr.TranslateError()
		} else {
			return r.serverError(err)
		}
	}
	chapter, err := r.Repository.Chapter.Get(int64(chapterNo), int64(bookId))
	if err != nil {
		switch {
		case errors.Is(err, utils.ErrorRecordsNotFound):
			return r.notFoundError(err)
		default:
			return r.serverError(err)
		}
	}
	if chapter.AuthorID != int64(userId) {
		return r.forbiddenError(utils.ErrorForbiddenResource)
	}
	if updateChapterPayload.Title != "" {
		chapter.Title = updateChapterPayload.Title
	}
	if updateChapterPayload.Description != "" {
		chapter.Description = updateChapterPayload.Description
	}
	err = r.Repository.Chapter.Update(chapter)
	if err != nil {
		r.serverError(err)
	}
	return c.JSON(http.StatusOK, Response[repositories.Chapter]{
		OK:   true,
		Data: *chapter,
	})
}

func (r Router) DeleteChapter(c echo.Context) error {
	chapterNo, err := strconv.Atoi(c.Param("chapterNo"))
	if err != nil {
		return r.badRequestError(err)
	}
	bookId, err := strconv.Atoi(c.Param("bookId"))
	if err != nil {
		return r.badRequestError(err)
	}
	userId, err := r.JwtService.RetrieveUserIdFromContext(c)
	if err != nil {
		return r.forbiddenError(err)
	}
	chapter, err := r.Repository.Chapter.Get(int64(chapterNo), int64(bookId))
	if err != nil {
		switch {
		case errors.Is(err, utils.ErrorRecordsNotFound):
			return r.notFoundError(err)
		default:
			return r.serverError(err)
		}
	}
	if chapter.AuthorID != int64(userId) {
		return r.forbiddenError(utils.ErrorForbiddenResource)
	}
	return c.JSON(http.StatusOK, Response[any]{
		OK: true,
	})
}
