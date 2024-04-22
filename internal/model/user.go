package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	FirstName string
	LastName  string
	Username  string `gorm:"unique"`
	Email     string `gorm:"unique"`
	Password  string
	// TODO: Implement rest of the fields
}

type UserQuery struct {
	ID       string `default:""` // It is just for querying, we can set it to a string
	Username string `default:""`
	Email    string `default:""`
}

type UserRegisterForm struct {
	FirstName string `validator:"required,min=3"`
	LastName  string `validator:"required,min=3"`
	Username  string `validator:"required,username,min=3,max=32"`
	Email     string `validator:"required,email"`
	Password  string `validator:"required,password,min=12,max=48"`
}
