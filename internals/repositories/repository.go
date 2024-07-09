package repositories

import (
	"strings"

	"github.com/jmoiron/sqlx"
)

type Repository struct {
	User    UserQueries
	Book    BookQueries
	Chapter ChapterQueries
	Content ContentQueries
}

func New(db *sqlx.DB) Repository {
	return Repository{
		User: UserRepository{
			DB: db,
		},
		Book: BookRepository{
			DB: db,
		},
		Chapter: ChapterRepository{
			DB: db,
		},
		Content: ContentRepository{
			DB: db,
		},
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
	TotalRecords int `json:"totalRecords"`
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
