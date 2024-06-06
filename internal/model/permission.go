package model

type Permission struct {
	ID   uint   `gorm:"primaryKey"`
	Name string `gorm:"unique"`
}

type UserPermission struct {
	UserID       uint `gorm:"Index"`
	PermissionID uint `gorm:"Index"`
}
