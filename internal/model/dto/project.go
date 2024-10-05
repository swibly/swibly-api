package dto

import (
	"time"

	"github.com/devkcud/arkhon-foundation/arkhon-api/pkg/utils"
	"gorm.io/gorm"
)

type ProjectCreation struct {
	Name        string `validate:"required,min=3,max=32"    json:"name"`
	Description string `validate:"omitempty,min=3,max=5000" json:"description"`

	Content any `validate:"omitempty" json:"content"`
	Budget  int `validate:"omitempty" json:"budget"`

	OwnerID uint `json:"-"` // Set using JWT

	Public bool `json:"-"` // Set in an URL query param
}

type ProjectUpdate struct {
	Name        *string `validate:"omitempty,min=3,max=32"    json:"name"`
	Description *string `validate:"omitempty,min=3,max=5000" json:"description"`

	Content *any `validate:"omitempty" json:"content"`
	Budget  *int `validate:"omitempty" json:"budget"`

	Published *bool `validate:"omitempty" json:"published"`
}

type AllowManage struct {
	Users    bool `gorm:"not null;default:false"`
	Metadata bool `gorm:"not null;default:false"`
}

type Allow struct {
	View    bool        `gorm:"not null;default:false"` // Will be ignored if project is public
	Edit    bool        `gorm:"not null;default:false"`
	Delete  bool        `gorm:"not null;default:false"`
	Publish bool        `gorm:"not null;default:false"`
	Share   bool        `gorm:"not null;default:false"` // Will be ignored if project is public
	Manage  AllowManage `gorm:"embedded;embeddedPrefix:manage_"`
}

type ProjectUserPermissions struct {
	ID             uint   `json:"id"`
	Username       string `json:"username"`
	ProfilePicture string `json:"profile_picture"`
	View           bool   `json:"allow_view"`
	Edit           bool   `json:"allow_edit"`
	Delete         bool   `json:"allow_delete"`
	Publish        bool   `json:"allow_publish"`
	Share          bool   `json:"allow_share"`
	ManageUsers    bool   `json:"allow_manage_users"`
	ManageMetadata bool   `json:"allow_manage_metadata"`
}

type ProjectInfoJSON struct {
	ID        uint           `json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at"`

	Name        string `json:"name"`
	Description string `json:"description"`

	Budget int `json:"budget"`

	IsPublic bool `json:"is_public"`

	OwnerID             uint   `json:"owner_id"`
	OwnerUsername       string `json:"owner_username"`
	OwnerProfilePicture string `json:"owner_profile_picture"`

	IsLiked          bool    `json:"is_liked"`
	IsDisliked       bool    `json:"is_disliked"`
	TotalLikes       int     `json:"total_likes"`
	TotalDislikes    int     `json:"total_dislikes"`
	LikeDislikeRatio float64 `json:"like_dislike_ratio"`

	AllowedUsers utils.JSON `gorm:"type:jsonb" json:"allowed_users"`
}

type ProjectInfo struct {
	ID        uint           `json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at"`

	Name        string `json:"name"`
	Description string `json:"description"`

	Budget int `json:"budget"`

	IsPublic bool `json:"is_public"`

	OwnerID             uint   `json:"owner_id"`
	OwnerUsername       string `json:"owner_username"`
	OwnerProfilePicture string `json:"owner_profile_picture"`

	IsLiked          bool    `json:"is_liked"`
	IsDisliked       bool    `json:"is_disliked"`
	TotalLikes       int     `json:"total_likes"`
	TotalDislikes    int     `json:"total_dislikes"`
	LikeDislikeRatio float64 `json:"like_dislike_ratio"`

	AllowedUsers []ProjectUserPermissions `json:"allowed_users"`
}

func (a Allow) IsEmpty() bool {
	return !a.View && !a.Edit && !a.Delete && !a.Publish && !a.Share && !a.Manage.Users && !a.Manage.Metadata
}
