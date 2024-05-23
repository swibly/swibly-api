package dto

type ProfileSearch struct {
	ID uint `json:"id"`

	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
	Username  string `json:"username"`

	Bio      string `json:"bio"`
	Verified bool   `json:"verified"`
}
