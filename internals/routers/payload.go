package router

import "gin_stuff/internals/data"

// -- REQUEST DATA
type LoginPayload struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type RegisterPayload struct {
	Username string `json:"username" validate:"required,min=6"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,strongPassword"`
}

// --- RESPONSE DATA
type LoginData struct {
	AccessToken string    `json:"accessToken"`
	User        data.User `json:"user"`
}

type RegisterData struct {
	AccessToken string    `json:"accessToken"`
	User        data.User `json:"user"`
}
