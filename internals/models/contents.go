package models

import (
	"context"
	"database/sql"
	"errors"
	"gin_stuff/internals/utils"
	"time"

	"github.com/jmoiron/sqlx"
)

type Content struct {
	ID          int64      `db:"id" json:"-"`
	ChapterID   int64      `db:"chapter_id" json:"-"`
	Chapter     *Chapter   `json:"chapter"`
	TextContent string     `db:"text_content" json:"textContent"`
	CreatedAt   *time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt   *time.Time `db:"updated_at" json:"updatedAt"`
	DeletedAt   *time.Time `db:"deleted_at" json:"deletedAt"`
}

type ContentRepostiory interface {
	Insert(*Content) error
	Get(int64) (*Content, error)
	Update(*Content) error
}

type ContentModel struct {
	DB *sqlx.DB
}

func (m ContentModel) Insert(content *Content) error {
	if len(content.TextContent) <= 0 {
		return utils.ErrorInvalidModel
	}
	statement := `
        INSERT INTO contents (chapter_id, text_content)
        VALUES ($1, $2)
        RETURNING id, created_at
    `

	args := []interface{}{content.ChapterID, content.TextContent}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	row := m.DB.QueryRowContext(ctx, statement, args...)
	return row.Scan(&content.ID, &content.CreatedAt)
}

// this basically only returns the latest content for the chapter
// different version of chapter content is stored in a different table
func (m ContentModel) Get(chapterID int64) (*Content, error) {
	if chapterID < 1 {
		return nil, utils.ErrorRecordsNotFound
	}
	statement := `
        SELECT 
            ct.id, ct.chapter_id, ct.text_content, ct.created_at, ct.updated_at,
            ch.id, ch.title, ch.chapter_no, ch.description, ch.created_at, ch.updated_at 
        FROM contents ct
        JOIN chapters ch ON ch.id = ct.chapter_id
        WHERE ch.id = $1 AND ch.deleted_at IS NULL
        LIMIT 1
    `
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	content := new(Content)
	content.Chapter = new(Chapter)
	row := m.DB.QueryRowContext(ctx, statement, chapterID)
	err := row.Scan(
		&content.ID, &content.ChapterID, &content.TextContent, &content.CreatedAt, &content.UpdatedAt,
		&content.Chapter.ID, &content.Chapter.Title, &content.Chapter.ChapterNO, &content.Chapter.Description, &content.Chapter.CreatedAt, &content.Chapter.UpdatedAt,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, utils.ErrorRecordsNotFound
		default:
			return nil, err
		}
	}
	return content, nil
}

func (m ContentModel) Update(content *Content) error {
    statement := `
        UPDATE contents
        SET text_content = $1
        WHERE id = $2
        RETURNING text_content, updated_at
    `
    ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
    defer cancel()

    args := []interface{}{content.ID, content.TextContent}
    row := m.DB.QueryRowContext(ctx, statement, args...)
    return  row.Scan(&content.TextContent, &content.UpdatedAt)
}
