package models

import "time"

type Content struct {
	ID          int64      `db:"id" json:"-"` // kinda dont need a id for content
	TextContent string     `db:"text_content" json:"textContent"`
	CreatedAt   *time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt   *time.Time `db:"updated_at" json:"updatedAt"`
	DeletedAt   *time.Time `db:"deleted_at" json:"deletedAt"`
}
