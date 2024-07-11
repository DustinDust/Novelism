package router

import (
	"gin_stuff/internals/data"
	"gin_stuff/internals/services"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

type Router struct {
	Queries       *data.Queries
	MailerService *services.MailerService
	JwtService    *services.JWTService
	LoggerService *services.LoggerService
}

func New(queries *data.Queries, mailerService *services.MailerService, loggerService *services.LoggerService) Router {
	return Router{
		Queries:       queries,
		MailerService: mailerService,
		LoggerService: loggerService,
		JwtService:    &services.JWTService{}, // recreate each router creation since it does not initiate any object instance
	}
}

type Filter struct {
	Page         int
	PageSize     int
	Sort         string
	SortSafeList []string
}

type Metadata struct {
	CurrentPage  int `json:"currentPage,omitempty"`
	PageSize     int `json:"pageSize,omitempty"`
	TotalRecords int `json:"totalRecords,omitempty"`
}

func CalculateMetadata(total, pageSize, page int) Metadata {
	if total == 0 {
		return Metadata{}
	}
	return Metadata{
		CurrentPage:  page,
		PageSize:     pageSize,
		TotalRecords: total,
	}
}

func (f Filter) SortColumn() string {
	for _, safeValue := range f.SortSafeList {
		if f.Sort == safeValue {
			return strings.TrimPrefix(f.Sort, "-")
		}
	}
	panic("unsafe sort parameter: " + f.Sort)
}

func (f Filter) SortDirection() string {
	if strings.HasPrefix(f.Sort, "-") {
		return "DESC"
	} else {
		return "ASC"
	}
}

func (f Filter) Limit() int {
	return f.PageSize
}

func (f Filter) Offset() int {
	return f.PageSize * (f.Page - 1)
}

type Response[T interface{}] struct {
	OK       bool     `json:"ok"`
	Data     T        `json:"data"`
	Metadata Metadata `json:"metadata,omitempty"`
}

// return http errors
func (r Router) serverError(err error) error {
	return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
}

func (r Router) notFoundError(err error) error {
	return echo.NewHTTPError(http.StatusNotFound, err.Error())
}

func (r Router) badRequestError(err error) error {
	return echo.NewHTTPError(http.StatusBadRequest, err.Error())
}

func (r Router) forbiddenError(err error) error {
	return echo.NewHTTPError(http.StatusForbidden, err.Error())
}

func (r Router) unauthorizedError(err error) error {
	return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
}
