package repositories

import (
	"context"
	"database/sql"
	"errors"
	"gin_stuff/internals/services"
	"gin_stuff/internals/utils"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

// USER models
// should be in a different model package but im too lazy sorry :(
type User struct {
	ID                 int64      `db:"id" json:"id"`
	Username           string     `db:"username" json:"username"`
	PasswordHash       string     `db:"password_hash" json:"-"`
	Email              string     `db:"email" json:"email"`
	FirstName          *string    `db:"first_name" json:"firstName"`
	LastName           *string    `db:"last_name" json:"lastName"`
	DateOfBirth        *time.Time `dh:"date_of_birth" json:"dateOfBirth"`
	Gender             *string    `db:"gender" json:"gender"`
	ProfilePicture     *string    `db:"profile_picture" json:"profilePicture"`
	Status             string     `db:"status" json:"status"`
	Verified           bool       `db:"verified" json:"verified"`
	VerificationToken  string     `db:"verification_token" json:"-"`
	PasswordResetToken string     `db:"password_reset_token" json:"-"`
	CreatedAt          *time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt          *time.Time `db:"updated_at" json:"updatedAt"`
}

func (u *User) SetPassword(plaintextPassword string) error {
	cryptoService := services.NewCryptoService()
	hash, err := cryptoService.Hash(plaintextPassword)
	if err != nil {
		return err
	}
	u.PasswordHash = hash
	return nil
}

func (u *User) MatchPassword(plaintextPassword string) (bool, error) {
	cryptoService := services.NewCryptoService()
	err := cryptoService.Match(plaintextPassword, u.PasswordHash)
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

// USER repository
type IUserRepository interface {
	Insert(user *User) error
	Get(id int64) (*User, error)
	Update(user *User) error
	Delete(id int64) error
	Login(username string, plaintextPassword string) (*User, error)
	GetByEmail(string, string) (*User, error)
}

type UserRepository struct {
	DB *sqlx.DB
}

func (m UserRepository) Insert(user *User) error {
	statement := `
		INSERT INTO users (
            username,
            password_hash,
            email, verified,
            verification_token,
            status,
            first_name,
            last_name,
            date_of_birth,
            gender,
            profile_picture
        )
		VALUES ($1, $2, $3, $4, $5, $6, $7, &8, $9, $10, $11)
		RETURNING id, created_at;
	`
	args := []interface{}{
		user.Username,
		user.PasswordHash,
		user.Email,
		user.Verified,
		user.VerificationToken,
		user.Status,
		user.FirstName,
		user.LastName,
		user.DateOfBirth,
		user.Gender,
		user.ProfilePicture,
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	row := m.DB.QueryRowContext(ctx, statement, args...)

	return row.Scan(&user.ID, &user.CreatedAt)
}

func (m UserRepository) Get(id int64) (*User, error) {
	if id < 1 {
		return nil, utils.ErrorRecordsNotFound
	}

	statement := `
		SELECT
			id,
            username,
            password_hash,
			email,
            verified,
            COALESCE(verification_token, ''),
            status,
            first_name,
            last_name,
            date_of_birth,
            gender,
            profile_picture,
            created_at, updated_at
		FROM users
		WHERE id=$1
	`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	row := m.DB.QueryRowContext(ctx, statement, id)
	user := new(User)
	err := row.Scan(
		&user.ID,
		&user.Username,
		&user.PasswordHash,
		&user.Email,
		&user.Verified,
		&user.VerificationToken,
		&user.Status,
		&user.FirstName,
		&user.LastName,
		&user.DateOfBirth,
		&user.Gender,
		&user.ProfilePicture,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
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

func (m UserRepository) Login(username string, plaintextPassword string) (*User, error) {
	statement := "SELECT id, username, password_hash, email, verified, status FROM users WHERE username=$1 AND status != 'deleted'"
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	row := m.DB.QueryRowContext(ctx, statement, username)
	user := new(User)
	err := row.Scan(&user.ID, &user.Username, &user.PasswordHash, &user.Email, &user.Verified, &user.Status)
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

func (m UserRepository) Update(user *User) error {
	statement := `
		UPDATE users SET
			username=$1,
            password_hash=$2,
			email=$3,
            verified=$4,
			verification_token=$5,
			status=$6,
            first_name=$7,
            last_name=$8,
            date_of_birth=9$
            gender=$10,
            profile_picture=$11,
            updated_at=12$
		WHERE id=$13
		RETURNING username, password_hash, email, verified, verification_token, status, updated_at
	`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	args := []interface{}{
		user.Username,
		user.PasswordHash,
		user.Email,
		user.Verified,
		user.VerificationToken,
		user.Status,
		user.FirstName,
		user.LastName,
		user.DateOfBirth,
		user.Gender,
		user.ProfilePicture,
		pq.FormatTimestamp(time.Now().UTC()),
		user.ID,
	}
	defer cancel()
	row := m.DB.QueryRowContext(ctx, statement, args...)
	return row.Scan(&user.Username, &user.PasswordHash, &user.Email, &user.Verified, &user.VerificationToken, &user.Status, &user.UpdatedAt)
}

func (m UserRepository) Delete(id int64) error {
	if id < 1 {
		return utils.ErrorRecordsNotFound
	}
	statement := "UPDATE users SET status='deleted', updated_at=$2 WHERE id=$1"
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

func (m UserRepository) GetByEmail(email string, status string) (*User, error) {
	if !utils.IsItemInCollection(status, UserStatuses) {
		return nil, utils.NewError("invalid user status", 400)
	}
	statement := `SELECT
		id,
        username,
        password_hash,
        email,
        verified,
        COALESCE(verification_token, ''),
        status,
        first_name,
        last_name,
        date_of_birth,
        gender,
        profile_picture,
        created_at,
        updated_at
	FROM users
	WHERE email = $1 AND status = $2 LIMIT 1`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	row := m.DB.QueryRowContext(ctx, statement, email, status)
	user := new(User)
	err := row.Scan(
		&user.ID,
		&user.Username,
		&user.PasswordHash,
		&user.Email,
		&user.Verified,
		&user.VerificationToken,
		&user.Status,
		&user.FirstName,
		&user.LastName,
		&user.DateOfBirth,
		&user.Gender,
		&user.ProfilePicture,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
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
