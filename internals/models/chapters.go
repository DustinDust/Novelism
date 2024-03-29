package models

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"gin_stuff/internals/utils"
	"math"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type Chapter struct {
	ID          int64      `db:"id" json:"id"`
	Book        *Book      `json:"-"`
	BookID      int64      `db:"book_id" json:"bookId"`
	Author      *User      `json:"-"`
	AuthorID    int64      `db:"author_id" json:"authorId"`
	ChapterNO   int64      `db:"chapter_no" json:"chapterNo"`
	Title       string     `db:"title" json:"title"`
	Content     *Content   `json:"content"`
	Description string     `db:"description" json:"description"`
	CreatedAt   *time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   *time.Time `db:"updated_at" json:"updated_at"`
	DeletedAt   *time.Time `db:"deleted_at" json:"deleted_at"`
}

type ChapterRepository interface {
	Insert(chapter *Chapter) error
	Get(chapterNo int64, bookId int64) (*Chapter, error)
	Update(chapter *Chapter) error
	Delete(id int64) error
	Find(bookId int64, title string, filter Filter) ([]*Chapter, Metadata, error)
	GetContent(chapterNo, bookId int64) (*Chapter, error)
	UpdateContent(chapter *Chapter, content string) error
}

type ChapterModel struct {
	DB *sqlx.DB
}

func (m ChapterModel) Find(bookId int64, title string, filter Filter) ([]*Chapter, Metadata, error) {
	if bookId < 1 {
		return nil, Metadata{}, utils.ErrorRecordsNotFound
	}
	statement := fmt.Sprintf(`
		SELECT count(*) OVER(), ch.id, ch.created_at, ch.updated_at, ch.deleted_at, ch.chapter_no, ch.title, ch.description, u.id, u.username, b.id
		FROM chapters ch
		JOIN users u ON u.id = ch.author_id
		JOIN books b ON b.id = ch.book_id
		WHERE b.id = $1 AND ch.deleted_at IS NULL
		AND (to_tsvector('simple', ch.title) @@ plainto_tsquery('simple', $2) OR $2 = '')
		ORDER BY %s %s, ch.chapter_no ASC
		LIMIT $3
		OFFSET $4
	`, filter.SortColumn(), filter.SortDirection())
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []interface{}{bookId, title, filter.Limit(), filter.Offset()}
	rows, err := m.DB.QueryContext(ctx, statement, args...)
	if err != nil {
		return nil, Metadata{}, err
	}
	chapters := []*Chapter{}
	totalRecords := 0
	for rows.Next() {
		var chapter Chapter
		chapter.Author = &User{}
		chapter.Book = &Book{}
		err := rows.Scan(
			&totalRecords,
			&chapter.ID,
			&chapter.CreatedAt,
			&chapter.UpdatedAt,
			&chapter.DeletedAt,
			&chapter.ChapterNO,
			&chapter.Title,
			&chapter.Description,
			&chapter.Author.ID,
			&chapter.Author.Username,
			&chapter.Book.ID,
		)
		if err != nil {
			return nil, Metadata{}, err
		}
		chapter.AuthorID = chapter.Author.ID
		chapter.BookID = chapter.Book.ID
		chapters = append(chapters, &chapter)
	}
	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err
	}
	return chapters, CalculateMetadata(totalRecords, filter.PageSize, filter.Page), nil
}

func (m ChapterModel) Insert(chapter *Chapter) error {
	if chapter.ChapterNO == 0 {
		chapters, _, err := m.Find(chapter.BookID, "", Filter{
			SortSafeList: []string{"-chapter_no"},
			Sort:         "-chapter_no",
			PageSize:     math.MaxInt,
			Page:         1,
		})
		if err != nil {
			return err
		}
		if len(chapters) > 0 {
			chapter.ChapterNO = chapters[0].ChapterNO + 1
		} else {
			chapter.ChapterNO = 1
		}
	}
	// this should create chapter only and the content will be added in later
	statement := `
		INSERT INTO chapters (book_id, author_id, chapter_no, title, description)
		VALUES($1, $2, $3, $4, $5)
		RETURNING id, created_at
	`

	args := []interface{}{chapter.BookID, chapter.AuthorID, chapter.ChapterNO, chapter.Title, chapter.Description}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	row := m.DB.QueryRowContext(ctx, statement, args...)
	return row.Scan(&chapter.ID, &chapter.CreatedAt)
}

// this get by the uniqe index, not the id. Personally I dont know what to do with it :(
func (m ChapterModel) Get(chapterNo int64, bookId int64) (*Chapter, error) {
	if bookId < 1 || chapterNo < 1 {
		return nil, utils.ErrorRecordsNotFound
	}

	statement := `
		SELECT ch.id, ch.title, ch.chapter_no,
		ch.description, ch.created_at, ch.updated_at,
		b.id, b.title, b.description,
		u.id, u.username, u.status, u.email
		FROM chapters ch
		JOIN books b ON b.id = ch.book_id
		JOIN users u ON u.id = ch.author_id
		WHERE ch.chapter_no = $1 AND b.id = $2 AND ch.deleted_at IS NULL
		LIMIT 1
	`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	chapter := new(Chapter)
	chapter.Author = new(User)
	chapter.Book = new(Book)
	row := m.DB.QueryRowContext(ctx, statement, chapterNo, bookId)
	err := row.Scan(
		&chapter.ID, &chapter.Title, &chapter.ChapterNO, &chapter.Description,
		&chapter.CreatedAt, &chapter.UpdatedAt, &chapter.Book.ID, &chapter.Book.Title,
		&chapter.Book.Description, &chapter.Author.ID, &chapter.Author.Username, &chapter.Author.Status, &chapter.Author.Email,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, utils.ErrorRecordsNotFound
		default:
			return nil, err
		}
	}
	// set some ID manually since we want to arshal these values
	chapter.AuthorID = chapter.Author.ID
	chapter.BookID = chapter.Book.ID
	return chapter, nil
}

func (m ChapterModel) Update(ch *Chapter) error {
	statement := `
		UPDATE chapters
		SET title=$2, description=$3, updated_at=$4
		WHERE id=$1
		RETURNING title, description, chapter_no, updated_at
	`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []interface{}{ch.ID, ch.Title, ch.Description, pq.FormatTimestamp(time.Now().UTC())}
	row := m.DB.QueryRowContext(ctx, statement, args...)
	return row.Scan(&ch.Title, &ch.Description, &ch.ChapterNO, &ch.UpdatedAt)
}

func (m ChapterModel) Delete(id int64) error {
	if id < 1 {
		return utils.ErrorRecordsNotFound
	}
	statement := `
		UPDATE chapter
		SET deleted_at=$2
		WHERE id=$1
	`
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

func (m ChapterModel) GetContent(chapterNo, bookId int64) (*Chapter, error) {
	chapter := new(Chapter)
	chapter.Content = new(Content)
	chapter.Book = new(Book)
	chapter.Author = new(User)

	statement := `
		SELECT ch.id, ch.chapter_no, ch.title, b.id, b.title, a.id, a.username,
		c.id, c.text_content, c.updated_at, c.created_at
		FROM chapter ch
		JOIN content c ON c.chapter_id = ch.id
		JOIN users a ON a.id = ch.author_id
		JOIN books b ON b.id = ch.book_id
		WHERE ch.chapter_no = $1
		AND b.id = $2
		AND deleted_at IS NULL
		LIMIT 1
	`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	args := []interface{}{chapterNo, bookId}
	row := m.DB.QueryRowContext(ctx, statement, args...)
	err := row.Scan(
		&chapter.ID,
		&chapter.ChapterNO,
		&chapter.Title,
		&chapter.Book.ID,
		&chapter.Book.Title,
		&chapter.Author.ID,
		&chapter.Author.Username,
		&chapter.Content.ID,
		&chapter.Content.TextContent,
		&chapter.Content.CreatedAt,
		&chapter.Content.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return chapter, nil
}

func (m ChapterModel) UpdateContent(chapter *Chapter, textContent string) error {
	// check if this chapter object already populate content
	if chapter.Content.ID == 0 {
		var chapterNo, bookId int64
		chapterNo = chapter.ChapterNO
		if chapter.BookID != 0 {
			bookId = chapter.BookID
		} else if chapter.Book.ID != 0 {
			bookId = chapter.Book.ID
		} else {
			return utils.ErrorInvalidModel
		}
		c, err := m.GetContent(chapterNo, bookId)
		if err != nil {
			return err
		}
		chapter.Content = c.Content
	}
	if chapter.Content.ID == 0 {
		// content doesnt exist
		createContentStatement := `
			INSERT INTO contents (text_content)
			VALUES ($1)
			RETURNING id, text_content, created_at
		`
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		row := m.DB.QueryRowContext(ctx, createContentStatement, textContent)
		return row.Scan(&chapter.Content.ID, &chapter.Content.TextContent, &chapter.Content.CreatedAt)
	} else {
		updateContentStatement := `
			UPDATE contents
			SET text_content=$2, updated_at=$3
			WHERE id=$1
			RETURNING text_content, updated_at
		`
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		row := m.DB.QueryRowContext(
			ctx, updateContentStatement,
			chapter.Content.ID,
			textContent,
			pq.FormatTimestamp(time.Now().UTC()),
		)
		return row.Scan(&chapter.Content.TextContent, &chapter.Content.UpdatedAt)
	}
}
