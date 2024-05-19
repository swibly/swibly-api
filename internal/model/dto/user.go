package dto

type UserRegister struct {
	FirstName string `validate:"required,min=3"                  json:"firstname"`
	LastName  string `validate:"required,min=3"                  json:"lastname"`
	Username  string `validate:"required,username,min=3,max=32"  json:"username"`
	Email     string `validate:"required,email"                  json:"email"`
	Password  string `validate:"required,password,min=12,max=48" json:"password"`
}

type UserLogin struct {
	Username string `validate:"omitempty,username,min=3,max=32" json:"username"`
	Email    string `validate:"omitempty,email"                 json:"email"`
	Password string `validate:"required,password,min=12,max=48" json:"password"`
}
