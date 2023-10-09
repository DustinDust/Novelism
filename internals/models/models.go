package models

import (
	"errors"

	"github.com/jmoiron/sqlx"
)

var (
	ErrorRecordNotFound = errors.New("record not found")
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
