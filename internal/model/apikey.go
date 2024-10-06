package model

import "time"

type APIKey struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time

	Key   string `json:"key"`
	Owner string `json:"owner"`

	EnabledKeyManage   int `json:"enabled_key_manage"   gorm:"default:-1"`
	EnabledAuth        int `json:"enabled_auth"         gorm:"default:-1"`
	EnabledSearch      int `json:"enabled_search"       gorm:"default:-1"`
	EnabledUserFetch   int `json:"enabled_user_fetch"   gorm:"default:-1"`
	EnabledUserActions int `json:"enabled_user_actions" gorm:"default:-1"`
	EnabledProjects    int `json:"enabled_projects"     gorm:"default:-1"`

	TimesUsed uint `json:"times_used" gorm:"default:0"`
	MaxUsage  uint `json:"max_usage"  gorm:"default:0"`
}
