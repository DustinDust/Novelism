package models

import (
	"github.com/jmoiron/sqlx"
)

type Models struct {
	User UserRepository
	Book BookRepository
}

func NewModels(db *sqlx.DB) Models {
	return Models{
		User: UserModel{
			DB: db,
		},
		Book: BookModel{
			DB: db,
		},
	}
}
