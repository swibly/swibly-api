package model

type APIKey struct {
	Key string `json:"key" gorm:"primarykey"`

	EnabledKeyManage   int `json:"enabled_key_manage"   gorm:"default:-1"` // Manage existing API keys
	EnabledAuth        int `json:"enabled_auth"         gorm:"default:-1"` // Register, Login, Update, Delete
	EnabledSearch      int `json:"enabled_search"       gorm:"default:-1"` // User
	EnabledUserFetch   int `json:"enabled_user_fetch"   gorm:"default:-1"` // Profile, Permissions, Followers, Following
	EnabledUserActions int `json:"enabled_user_actions" gorm:"default:-1"` // Follow, Unfollow

	TimesUsed uint `json:"times_used" gorm:"default:0"`
	MaxUsage  uint `json:"max_usage"  gorm:"default:0"`
}
