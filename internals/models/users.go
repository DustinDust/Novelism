package models

import (
	"context"
	"database/sql"
	"errors"
	"gin_stuff/internals/utils"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

// user models
type User struct {
	ID           int64      `db:"id" json:"id"`
	Username     string     `db:"username" json:"username"`
	PasswordHash string     `db:"password_hash" json:"-"`
	Email        string     `db:"email" json:"email"`
	Status       string     `db:"status" json:"-"`
	CreatedAt    *time.Time `db:"created_at" json:"created_at"`
	UpdatedAt    *time.Time `db:"updated_at" json:"updated_at"`
}

type UserRepository interface {
	Insert(user *User) error
	Get(id int64) (*User, error)
	Update(user *User) error
	Delete(id int64) error
	Login(username string, plaintextPassword string) (*User, error)
}

type UserModel struct {
	DB *sqlx.DB
}

func (m UserModel) Insert(user *User) error {
	statement := `
		INSERT INTO users (username, password_hash, email, status)
		VALUES ($1, $2, $3, &4)
		RETURNING id, created_at;
	`
	args := []interface{}{user.Username, user.PasswordHash, user.Email, user.Status}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	row := m.DB.QueryRowContext(ctx, statement, args...)

	return row.Scan(&user.ID, &user.CreatedAt)
}

func (m UserModel) Get(id int64) (*User, error) {
	if id < 1 {
		return nil, utils.ErrorRecordsNotFound
	}

	statement := "SELECT id, username, password_hash, email, created_at FROM users WHERE id=$1"
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	row := m.DB.QueryRowContext(ctx, statement, id)
	user := new(User)
	err := row.Scan(&user.ID, &user.Username, &user.PasswordHash, &user.Email, &user.CreatedAt)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, utils.ErrorRecordsNotFound
		default:
			return nil, err
		}
	}
	return user, nil
}

func (m UserModel) Login(username string, plaintextPassword string) (*User, error) {
	statement := "SELECT id, username, password_hash, email FROM users WHERE username=$1 AND status != 'deleted'"
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	row := m.DB.QueryRowContext(ctx, statement, username)
	user := new(User)
	err := row.Scan(&user.ID, &user.Username, &user.PasswordHash, &user.Email)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, utils.ErrorInvalidCredentials
		default:
			return nil, err
		}
	}
	match, err := user.MatchPassword(plaintextPassword)
	if err != nil {
		return nil, err
	}
	if !match {
		return nil, utils.ErrorInvalidCredentials
	}
	return user, nil
}

func (m UserModel) Update(user *User) error {
	statement := `
		UPDATE users
		SET username=$1, password_hash=$2, email=$3, status=$4, updated_at=$5
		WHERE id=$5
		RETURNING username, password_hash, email, status, updated_at
	`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	args := []interface{}{user.Username, user.PasswordHash, user.Email, user.Status, pq.FormatTimestamp(time.Now().UTC())}
	defer cancel()
	row := m.DB.QueryRowContext(ctx, statement, args...)
	return row.Scan(&user.Username, &user.PasswordHash, &user.Email, &user.UpdatedAt)
}

func (m UserModel) Delete(id int64) error {
	if id < 1 {
		return utils.ErrorRecordsNotFound
	}
	statement := "UPDATE users SET status='deleted' WHERE id=$1"
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := m.DB.ExecContext(ctx, statement, id)
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

func (u *User) SetPassword(plaintextPassword string) error {
	hash, err := utils.Hash(plaintextPassword)
	if err != nil {
		return err
	}
	u.PasswordHash = hash
	return nil
}

func (u *User) MatchPassword(plaintextPassword string) (bool, error) {
	err := utils.Match(plaintextPassword, u.PasswordHash)
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil
		default:
			return false, err
		}
	}
	return true, nil
}
