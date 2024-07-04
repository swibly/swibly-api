package model

type APIKey struct {
	Key string `gorm:"primarykey"`

	EnabledKeyManage   bool `gorm:"default:false"` // Manage existing API keys
	EnabledAuth        bool `gorm:"default:false"` // Register, Login, Update, Delete
	EnabledSearch      bool `gorm:"default:false"` // User
	EnabledUserFetch   bool `gorm:"default:false"` // Profile, Permissions, Followers, Following
	EnabledUserActions bool `gorm:"default:false"` // Follow, Unfollow
}
