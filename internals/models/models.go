package models

import (
	"github.com/jmoiron/sqlx"
)

type Models struct {
	User UserRepository
}

func NewModels(db *sqlx.DB) Models {
	return Models{
		User: UserModel{
			DB: db,
		},
	}
}
