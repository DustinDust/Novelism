package data

import (
	"encoding/json"

	"github.com/jackc/pgx/v5/pgtype"
)

// Implement custom marshaler for types generated by the `sqlc genreate` command

func (u User) MarshalJSON() ([]byte, error) {
	type JSONUser struct {
		ID                 int32            `json:"id"`
		Username           string           `json:"username"`
		PasswordHash       string           `json:"-"`
		Email              string           `json:"email"`
		CreatedAt          pgtype.Timestamp `json:"created_at"`
		UpdatedAt          pgtype.Timestamp `json:"updated_at"`
		Status             UserStatus       `json:"status"`
		Verified           pgtype.Bool      `json:"verified"`
		VerificationToken  pgtype.Text      `json:"-"`
		PasswordResetToken pgtype.Text      `json:"-"`
		FirstName          pgtype.Text      `json:"first_name"`
		LastName           pgtype.Text      `json:"last_name"`
		DateOfBirth        pgtype.Date      `json:"date_of_birth"`
		Gender             pgtype.Text      `json:"gender"`
		ProfilePicture     pgtype.Text      `json:"profile_picture"`
	}

	return json.Marshal(JSONUser{
		ID:             u.ID,
		Username:       u.Username,
		Email:          u.Email,
		CreatedAt:      u.CreatedAt,
		UpdatedAt:      u.UpdatedAt,
		Status:         u.Status.UserStatus,
		Verified:       u.Verified,
		FirstName:      u.FirstName,
		LastName:       u.LastName,
		DateOfBirth:    u.DateOfBirth,
		Gender:         u.Gender,
		ProfilePicture: u.ProfilePicture,
	})
}

func (status NullUserStatus) MarshalJSON() ([]byte, error) {
	if !status.Valid {
		return nil, nil
	}
	return []byte(status.UserStatus), nil
}

func (v NullVisibility) MarshalJSON() ([]byte, error) {
	if !v.Valid {
		return nil, nil
	}
	return []byte(v.Visibility), nil
}
