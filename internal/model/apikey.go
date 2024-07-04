package model

type APIKey struct {
	Key string `gorm:"primarykey"`

	EnabledKeyManage   int `gorm:"default:-1"` // Manage existing API keys
	EnabledAuth        int `gorm:"default:-1"` // Register, Login, Update, Delete
	EnabledSearch      int `gorm:"default:-1"` // User
	EnabledUserFetch   int `gorm:"default:-1"` // Profile, Permissions, Followers, Following
	EnabledUserActions int `gorm:"default:-1"` // Follow, Unfollow
}
