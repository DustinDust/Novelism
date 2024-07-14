package router

import (
	"gin_stuff/internals/data"
	"gin_stuff/internals/utils"
	"net/http"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/labstack/echo/v4"
)

func (r Router) GetBooks(c echo.Context) error {
	filter := Filter{}
	if err := c.Bind(&filter); err != nil {
		return r.badRequestError(err)
	}
	user, ok := c.Get("user").(data.User)
	if !ok {
		return r.forbiddenError(utils.ErrorForbiddenResource())
	}
	totalBooks, err := r.queries.CountBooksByUserId(c.Request().Context(), pgtype.Int4{Int32: user.ID, Valid: true})
	if err != nil {
		return r.serverError(err)
	}
	books, err := r.queries.FindBooksByUserId(c.Request().Context(), data.FindBooksByUserIdParams{
		UserID: pgtype.Int4{Int32: user.ID, Valid: true},
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
