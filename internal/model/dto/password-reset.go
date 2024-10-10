package dto

type RequestPasswordReset struct {
	Email string `validate:"required,email" json:"email"`
}

type PasswordReset struct {
	Password string `validate:"required,password" json:"password"`
}
