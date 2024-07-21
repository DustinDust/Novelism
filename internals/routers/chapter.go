package router

import (
	"errors"
	"gin_stuff/internals/data"
	"net/http"
	"strconv"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/labstack/echo/v4"
)

func (r Router) GetChapters(c echo.Context) error {
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
