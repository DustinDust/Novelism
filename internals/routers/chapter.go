package router

import (
	"context"
	"errors"
	"gin_stuff/internals/data"
	"gin_stuff/internals/utils"
	"net/http"
	"strconv"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/labstack/echo/v4"
)

func (r Router) GetChaptersByBook(c echo.Context) error {
	bookID, err := strconv.Atoi(c.Param("bookId"))
	if err != nil {
		return r.badRequestError(err)
	}
	book, err := r.queries.GetBookById(c.Request().Context(), int32(bookID))
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return r.notFoundError(err)
		default:
			return r.serverError(err)
		}
	}
	chapters, err := r.queries.FindChaptersByBookId(c.Request().Context(), pgtype.Int4{Int32: int32(bookID), Valid: true})
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return c.JSON(http.StatusOK, Response[GetChaptersData]{
				OK: true,
				Data: GetChaptersData{
					Book:     book,
					Chapters: []data.Chapter{},
				},
			})
		default:
			return r.serverError(err)
		}
	}
	return c.JSON(http.StatusOK, Response[GetChaptersData]{
		OK: true,
		Data: GetChaptersData{
			Book:     book,
			Chapters: chapters,
		},
	})
}

func (r Router) CreateChapter(c echo.Context) error {
	bookID, err := strconv.Atoi(c.Param("bookId"))
	if err != nil {
		return r.badRequestError(err)
	}
	user, ok := c.Get("user").(data.User)
	if !ok {
		return r.unauthorizedError(utils.ErrUnauthorized())
	}
	book, err := r.queries.GetBookById(c.Request().Context(), int32(bookID))
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return r.notFoundError(err)
		default:
			return r.serverError(err)
		}
	}
	if book.UserID.Int32 != user.ID {
		return r.unauthorizedError(utils.ErrorForbiddenResource("book", "chapters"))
	}
	payload := CreateChapterPayload{}
	if err := r.bindAndValidatePayload(c, &payload); err != nil {
		return r.badRequestError(err)
	}
	tx, err := r.db.BeginTx(c.Request().Context(), pgx.TxOptions{})
	if err != nil {
		return r.serverError(err)
	}
	defer tx.Rollback(context.Background())

	chapter, err := r.queries.WithTx(tx).InsertChapter(c.Request().Context(), data.InsertChapterParams{
		BookID:      pgtype.Int4{Int32: book.ID, Valid: true},
		AuthorID:    pgtype.Int4{Int32: user.ID, Valid: true},
		Title:       pgtype.Text{String: payload.Title, Valid: true},
		Description: pgtype.Text{String: payload.Description, Valid: true},
	})
	if err != nil {
		tx.Rollback(c.Request().Context())
		return r.serverError(err)
	}
	content, err := r.queries.InsertContentToChapter(c.Request().Context(), data.InsertContentToChapterParams{
		ChapterID:   pgtype.Int4{Int32: chapter.ID, Valid: true},
		TextContent: pgtype.Text{String: "", Valid: true},
		Status:      data.NullContentStatus{ContentStatus: data.ContentStatusDraft, Valid: true},
	})
	if err != nil {
		tx.Rollback(c.Request().Context())
		return r.serverError(err)
	}
	tx.Commit(c.Request().Context())

	return c.JSON(http.StatusOK, Response[CreateChapterData]{
		OK:   true,
		Data: CreateChapterData{Chapter: chapter, Content: []data.Content{content}},
	})
}

func (r Router) ChapterDetail(c echo.Context) error {
	return echo.ErrNotImplemented
}
