package dto

import "time"

type RequestPasswordReset struct {
	Email string `validate:"required,email" json:"email"`
}

type PasswordReset struct {
	Password string `validate:"required,password" json:"password"`
}

type PasswordResetInfo struct {
	FirstName      string    `json:"firstname"`
	LastName       string    `json:"lastname"`
	Username       string    `json:"username"`
	ProfilePicture string    `json:"profile_picture"`
	Lang           string    `json:"lang"`
	ExpiresAt      time.Time `json:"expires_at"`
}
