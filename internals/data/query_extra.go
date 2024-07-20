package data

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

// manual method that i can't figure out a way to generate nicely

// Save book to DB
func (q *Queries) SaveBook(ctx context.Context, book *Book) error {
	statement := `
		UPDATE books
		SET title=$1, description=$2, updated_at=$3
		WHERE id=$4
		RETURNING title, description, updated_at
	`
	args := []interface{}{book.Title, book.Description, pgtype.Timestamp{Time: time.Now(), Valid: true}, book.ID}
	row := q.db.QueryRow(ctx, statement, args...)
	return row.Scan(&book.Title, &book.Description, &book.UpdatedAt)
}
