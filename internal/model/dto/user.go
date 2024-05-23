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

type UserUpdate struct {
	FirstName string `validate:"omitempty,min=3"                 json:"firstname"`
	LastName  string `validate:"omitempty,min=3"                 json:"lastname"`
	Username  string `validate:"omitempty,username,min=3,max=32" json:"username"`

	Bio      string `validate:"omitempty,max=480" json:"bio"`
	Verified bool   `validate:"omitempty"         json:"verified"`

	// NOTE: XP and Arkhoins doesn't make sense to update here

	Email    string `validate:"omitempty,email"                  json:"email"`
	Password string `validate:"omitempty,password,min=12,max=48" json:"password"`

	// NOTE: Notification and Show structs will be in another method/route to update
	//       due to the nature of structs being a pain in the ass to work with :/
}
