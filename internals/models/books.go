package models

import (
	"context"
	"database/sql"
	"errors"
	"gin_stuff/internals/utils"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type Book struct {
	ID          int64      `db:"id" json:"id"`
	UserID      int64      `db:"user_id" json:"-"`
	User        *User      `json:"user"`
	Title       string     `db:"title" json:"title"`
	Description string     `db:"description" json:"description"`
	CreatedAt   *time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   *time.Time `db:"updated_at" json:"updated_at"`
}

type BookFilter struct {
	UserID      int64
	Title       string
	Description string
}

type BookRepository interface {
	Insert(book *Book) error
	Get(id int64) (*Book, error)
	Update(book *Book) error
	Delete(id int64) error
	Find(filter BookFilter) ([]*Book, error)
}

type BookModel struct {
	DB *sqlx.DB
}

func (m BookModel) Insert(book *Book) error {
	statement := `
		INSERT INTO books (title, description, user_id)
		VALUES ($1, $2, $3)
		RETURNING id, created_at, user_id
	`
	args := []interface{}{book.Title, book.Description, book.User.ID}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	row := m.DB.QueryRowContext(ctx, statement, args...)
	return row.Scan(&book.ID, &book.CreatedAt, &book.UserID)
}

func (m BookModel) Get(id int64) (*Book, error) {
	if id < 1 {
		return nil, utils.ErrorRecordsNotFound
	}
	statement := `
		SELECT b.id, b.title, b.description, b.created_at, b.updated_at, u.id, u.username, u.email
		FROM books b
		JOIN users u
		ON b.user_id = u.id
		WHERE b.id=$1
		LIMIT 1
	`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	book := new(Book)
	book.User = new(User)

	row := m.DB.QueryRowContext(ctx, statement, id)
	err := row.Scan(
		&book.ID,
		&book.Title,
		&book.Description,
		&book.CreatedAt,
		&book.UpdatedAt,
		&book.User.ID,
		&book.User.Username,
		&book.User.Email,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, utils.ErrorRecordsNotFound
		default:
			return nil, err
		}
	}
	book.UserID = book.User.ID // set UesrID since scanning does not automatically do this
	return book, nil
}

func (m BookModel) Update(b *Book) error {
	statement := `
		UPDATE books
		SET title=$1, description=$2, updated_at=$3
		WHERE id=$4
		RETURNING title, description, updated_at
	`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	args := []interface{}{b.Title, b.Description, pq.FormatTimestamp(time.Now().UTC()), b.ID}
	row := m.DB.QueryRowContext(ctx, statement, args...)
	return row.Scan(&b.Title, &b.Description, &b.UpdatedAt)
}

func (m BookModel) Delete(id int64) error {
	if id < 1 {
		return utils.ErrorRecordsNotFound
	}
	statement := "UPDATE books SET deleted_at=$2 WHERE id=$1;"
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := m.DB.ExecContext(ctx, statement, id, pq.FormatTimestamp(time.Now().UTC()))
	if err != nil {
		return err
	}
	rowAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowAffected == 0 {
		return utils.ErrorRecordsNotFound
	}
	return nil
}

func (m BookModel) Find(filter BookFilter) ([]*Book, error) {
	return []*Book{}, nil
}
