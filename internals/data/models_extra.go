package data

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

func (b *Book) PatchAttrs(attrs map[string]any) {
	if t, ok := attrs["title"].(string); ok && t != "" {
		b.Title = pgtype.Text{String: t, Valid: true}
	}
	if d, ok := attrs["description"].(string); ok && d != "" {
		b.Description = pgtype.Text{String: d, Valid: true}
	}
	if d, ok := attrs["cover"].(string); ok && d != "" {
		b.Cover = pgtype.Text{String: d, Valid: true}
	}
}

func (u *User) PatchAttrs(attrs map[string]any) {
	if lastName, ok := attrs["lastName"].(string); ok && lastName != "" {
		u.LastName = pgtype.Text{String: lastName, Valid: true}
	}
	if firstName, ok := attrs["firstName"].(string); ok && firstName != "" {
		u.LastName = pgtype.Text{String: firstName, Valid: true}
	}
	if gender, ok := attrs["gender"].(string); ok && gender != "" {
		u.Gender = pgtype.Text{String: gender, Valid: true}
	}
	if profilePicture, ok := attrs["profilePicture"].(string); ok && profilePicture != "" {
		u.Gender = pgtype.Text{String: profilePicture, Valid: true}
	}
	if dateOfBirth, ok := attrs["dateOfBirth"].(time.Time); ok {
		u.DateOfBirth = pgtype.Date{Time: dateOfBirth, Valid: true}
	}
}
