package router

import (
	"fmt"
	"gin_stuff/internals/data"
	"gin_stuff/internals/services"
	"gin_stuff/internals/utils"
	"net/http"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
)

type Router struct {
	db      *pgx.Conn
	queries *data.Queries

	mailer    *services.MailerService
	jwt       *services.JWTService
	validator *utils.Validator
}

func New(dbtx *pgx.Conn) (*Router, error) {
	mailer, err := services.NewMailerService(services.MailerSMTPConfig{
		Host:     viper.GetString("mailer.host"),
		Port:     viper.GetInt64("mailer.port"),
		Login:    viper.GetString("mailer.login"),
		Password: viper.GetString("mailer.password"),
		Timeout:  viper.GetDuration("mailer.timeout"),
	})
	if err != nil {
		return nil, err
	}
	return &Router{
		db:        dbtx,
		queries:   data.New(dbtx),
		mailer:    mailer,
		jwt:       &services.JWTService{}, // recreate each router creation since it does not initiate any object instance
		validator: utils.NewValidator(),
	}, nil
}

type Filter struct {
	Page         int    `query:"page"`
	PageSize     int    `query:"pageSize"`
	Sort         string `query:"sort"`
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
		return "desc"
	} else {
		return "asc"
	}
}

func (f Filter) SortString() string {
	return fmt.Sprintf("%s %s", f.SortColumn(), f.SortDirection())
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
