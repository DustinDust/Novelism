package router

import (
	"errors"
	"gin_stuff/internals/models"
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

type UpdateChapterContentPayload struct {
	TextContent string `json:"textContent" validate:"required"`
}

func getChapterNoAndBookId(c echo.Context) (int, int, error) {
	chapterNo, err := strconv.Atoi(c.Param("chapterNo"))
	if err != nil {
		return 0, 0, err
	}
	bookId, err := strconv.Atoi(c.Param("bookId"))
	if err != nil {
		return 0, 0, err
	}
	return chapterNo, bookId, nil
}

func (r Router) CreateChapter(c echo.Context) error {
	validate := utils.NewValidator()
	createChapterPayload := new(CreateChapterPayload)
	idStr := c.Param("bookId")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return r.badRequestError(err)
	}
	userId, err := utils.JWT.RetreiveUserIdFromContext(c)
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
	chapter := models.Chapter{
		AuthorID:    book.UserID,
		BookID:      book.ID,
		Title:       createChapterPayload.Title,
		Description: createChapterPayload.Description,
	}
	err = r.Model.Chapter.Insert(&chapter)
	if err != nil {
		return r.serverError(err)
	}
	return c.JSON(http.StatusCreated, echo.Map{
		"data": chapter,
		"ok":   true,
	})
}

// get all chapters from 1 book (with optional filter)
func (r Router) FindChapters(c echo.Context) error {
	bookIdStr := c.Param("bookId")
	bookId, err := strconv.Atoi(bookIdStr)
	if err != nil {
		return r.badRequestError(err)
	}
	_, err = r.Model.Book.Get(int64(bookId))
	if err != nil {
		switch {
		case errors.Is(err, utils.ErrorRecordsNotFound):
			return r.notFoundError(err)
		default:
			return r.serverError(err)
		}
	}
	var title string
	filter := models.Filter{
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

	chapters, metadata, err := r.Model.Chapter.Find(int64(bookId), title, filter)
	if err != nil {
		return r.serverError(err)
	}
	return c.JSON(http.StatusOK, echo.Map{
		"metadata": metadata,
		"data": echo.Map{
			"chapters": chapters,
		},
		"ok": true,
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
	userId, err := utils.JWT.RetreiveUserIdFromContext(c)
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
	chapter, err := r.Model.Chapter.Get(int64(chapterNo), int64(bookId))
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
	err = r.Model.Chapter.Update(chapter)
	if err != nil {
		r.serverError(err)
	}
	return c.JSON(http.StatusOK, echo.Map{
		"ok":   true,
		"data": chapter,
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
	userId, err := utils.JWT.RetreiveUserIdFromContext(c)
	if err != nil {
		return r.forbiddenError(err)
	}
	chapter, err := r.Model.Chapter.Get(int64(chapterNo), int64(bookId))
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
	if err != nil {
		return r.serverError(err)
	}
	return c.JSON(http.StatusOK, echo.Map{"ok": true})
}

func (r Router) GetChapterContent(c echo.Context) error {
	chapterNo, err := strconv.Atoi(c.Param("chapterNo"))
	if err != nil {
		return r.badRequestError(err)
	}
	bookId, err := strconv.Atoi(c.Param("bookId"))
	if err != nil {
		return r.badRequestError(err)
	}
	// userId, err := utils.RetreiveUserIdFromContext(c)
	// if err != nil {
	// 	return r.forbiddenError(err)
	// }
	chapter, err := r.Model.Chapter.GetContent(int64(chapterNo), int64(bookId))
	if err != nil {
		switch {
		case errors.Is(err, utils.ErrorRecordsNotFound):
			return r.notFoundError(err)
		default:
			return r.serverError(err)
		}
	}
	return c.JSON(http.StatusOK, echo.Map{
		"ok":   true,
		"data": chapter,
	})
}

func (r Router) UpdateChapterContent(c echo.Context) error {
	chapterNo, bookId, err := getChapterNoAndBookId(c)
	if err != nil {
		return r.badRequestError(err)
	}
	userId, err := utils.JWT.RetreiveUserIdFromContext(c)
	if err != nil {
		return r.forbiddenError(err)
	}
	chapter, err := r.Model.Chapter.GetContent(int64(chapterNo), int64(bookId))
	if err != nil {
		return r.serverError(err)
	}
	if int64(userId) != chapter.Author.ID {
		return r.forbiddenError(utils.ErrorForbiddenResource)
	}
	validate := utils.NewValidator()
	payload := new(UpdateChapterContentPayload)
	if err := c.Bind(&payload); err != nil {
		r.badRequestError(err)
	}
	if err := validate.ValidateStruct(payload); err != nil {
		if verr, ok := err.(*utils.StructValidationErrors); ok {
			return verr.TranslateError()
		} else {
			return r.serverError(err)
		}
	}
	if err := r.Model.Chapter.UpdateContent(chapter, payload.TextContent); err != nil {
		r.serverError(err)
	}
	return c.JSON(http.StatusOK, echo.Map{
		"ok":   true,
		"data": chapter,
	})
}
