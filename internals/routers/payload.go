package router

import "gin_stuff/internals/data"

// -- REQUEST DATA
type SignInPayload struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type SignUpPayload struct {
	Username string `json:"username" validate:"required,min=6"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,strongPassword"`
}

type InsertBookPayload struct {
	Title       string `json:"title" validate:"required"`
	Description string `json:"description" validate:"required"`
}

type UpdateBookPayload struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type CreateChapterPayload struct {
	Title       string `json:"title" validate:"required"`
	Description string `json:"description" validate:"required"`
}

// --- RESPONSE DATA
type SignInData struct {
	AccessToken string    `json:"accessToken,omitempty"`
	User        data.User `json:"user"`
}

type SignUpData struct {
	AccessToken string    `json:"accessToken,omitempty"`
	User        data.User `json:"user"`
}

type GetChaptersData struct {
	Book     data.Book      `json:"book"`
	Chapters []data.Chapter `json:"chapters"`
}

type CreateChapterData struct {
	Chapter data.Chapter   `json:"chapter"`
	Content []data.Content `json:"content"`
}
