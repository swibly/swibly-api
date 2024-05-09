package model

import "gorm.io/gorm"

type Role struct {
	gorm.Model
	Name  string
	Users []User `gorm:"many2many:user_roles"`
}
