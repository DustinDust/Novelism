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

func (r Router) BrowseBooks(c echo.Context) error {
	filter := Filter{}
	if err := c.Bind(&filter); err != nil {
		return r.badRequestError(err)
	}
	totalBooks, err := r.queries.CountBrowsableBooks(c.Request().Context())
	if err != nil {
		return r.serverError(err)
	}
	books, err := r.queries.BrowseBooks(c.Request().Context(), data.BrowseBooksParams{
		Limit:  int32(filter.Limit()),
		Offset: int32(filter.Offset()),
	})
	if err != nil {
		return r.serverError(err)
	}
	return c.JSON(http.StatusOK, Response[[]data.Book]{
		OK:       true,
		Data:     books,
		Metadata: CalculateMetadata(int(totalBooks), filter.PageSize, filter.Page),
	})
}

func (r Router) CreateBook(c echo.Context) error {
	payload := InsertBookPayload{}
	if err := r.bindAndValidatePayload(c, &payload); err != nil {
		return err
	}
	user, ok := c.Get("user").(data.User)
	if !ok {
		return r.forbiddenError(utils.ErrorForbiddenResource())
	}

	tx, err := r.db.BeginTx(c.Request().Context(), pgx.TxOptions{})
	if err != nil {
		return r.serverError(err)
	}
	defer tx.Rollback(context.Background())

	book, err := r.queries.WithTx(tx).InsertBook(c.Request().Context(), data.InsertBookParams{
		UserID:      pgtype.Int4{Int32: user.ID, Valid: true},
		Title:       pgtype.Text{String: payload.Title, Valid: true},
		Description: pgtype.Text{String: payload.Description, Valid: payload.Description != ""},
	})
	if err != nil {
		return r.serverError(err)
	}

	if err := tx.Commit(c.Request().Context()); err != nil {
		return r.serverError(err)
	}

	return c.JSON(http.StatusCreated, Response[data.Book]{
		OK:   true,
		Data: book,
	})
}

func (r Router) UpdateBook(c echo.Context) error {
	bookId, err := strconv.Atoi(c.Param("bookId"))
	if err != nil {
		return r.badRequestError(err)
	}
	payload := UpdateBookPayload{}
	if err := r.bindAndValidatePayload(c, &payload); err != nil {
		return err
	}
	user, ok := c.Get("user").(data.User)
	if !ok {
		return r.forbiddenError(utils.ErrorForbiddenResource())
	}
	book, err := r.queries.GetBookById(c.Request().Context(), int32(bookId))
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return r.notFoundError(err)
		default:
			return r.serverError(err)
		}
	}
	if !book.UserID.Valid || book.UserID.Int32 != user.ID {
		return r.forbiddenError(utils.ErrorForbiddenResource())
	}

	tx, err := r.db.BeginTx(c.Request().Context(), pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer tx.Rollback(context.Background())

	book.PatchAttrs(map[string]any{
		"description": payload.Description,
		"title":       payload.Title,
	})
	if err := r.queries.WithTx(tx).SaveBook(c.Request().Context(), &book); err != nil {
		return r.serverError(err)
	}

	if err := tx.Commit(c.Request().Context()); err != nil {
		return r.serverError(err)
	}

	return c.JSON(http.StatusOK, Response[data.Book]{
		OK:   true,
		Data: book,
	})
}

func (r Router) DeleteBook(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("bookId"))
	if err != nil {
		return r.badRequestError(err)
	}
	user, ok := c.Get("user").(data.User)
	if !ok {
		return r.forbiddenError(utils.ErrorForbiddenResource())
	}
	book, err := r.queries.GetBookById(c.Request().Context(), int32(id))
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return r.notFoundError(err)
		default:
			return r.serverError(err)
		}
	}
	if book.UserID.Int32 != user.ID {
		return r.forbiddenError(utils.ErrorForbiddenResource())
	}
	err = r.queries.DeleteBook(c.Request().Context(), book.ID)
	if err != nil {
		return r.serverError(err)
	}
	return c.JSON(http.StatusOK, Response[any]{
		OK:   true,
		Data: book.ID,
	})
}
