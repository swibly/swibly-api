package auth

type UserBodyRegister struct {
	Fullname string `json:"fullname"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserBodyLogin struct {
	Username string `json:"username" default:""`
	Email    string `json:"email" default:""`
	Password string `json:"password" default:""`
}
