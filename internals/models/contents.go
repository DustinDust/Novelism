package models

import "time"

type Content struct {
	ID          int64      `db:"id" json:"-"` // kinda dont need a id for content
	TextContent string     `db:"text_content" json:"textContent"`
	CreatedAt   *time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   *time.Time `db:"updated_at" json:"updated_at"`
	DeletedAt   *time.Time `db:"deleted_at" json:"deleted_at"`
}
