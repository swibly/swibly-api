package dto

import (
	"mime/multipart"
	"time"

	"github.com/swibly/swibly-api/pkg/utils"
	"gorm.io/gorm"
)

type ProjectCreation struct {
	Name        string `validate:"required,min=3,max=32"    form:"name"`
	Description string `validate:"omitempty,min=3,max=5000" form:"description"`

	BannerImage *multipart.FileHeader `validate:"omitempty" form:"banner"`

	Content any `validate:"omitempty" form:"content"`
	Budget  int `validate:"omitempty,max=1000000000000" form:"budget"`

	Width  int `validate:"omitempty,min=1,max=1000" form:"width"`
	Height int `validate:"omitempty,min=1,max=1000" form:"height"`

	OwnerID uint `form:"-"` // Set using JWT

	Public bool `form:"-"` // Set in an URL query param

	Fork *uint `form:"-"` // Set in code during the clone procedure
}

type ProjectUpdate struct {
	Name        *string `validate:"omitempty,min=3,max=32" form:"name"`
	Description *string `validate:"omitempty,max=5000"     form:"description"`

	BannerImage *multipart.FileHeader `validate:"omitempty" form:"banner"`

	Content any  `validate:"omitempty" form:"-"` // Set in code
	Budget  *int `validate:"omitempty,max=1000000000000" form:"budget"`

	Width  *int `validate:"omitempty,min=1,max=1000" form:"width"`
	Height *int `validate:"omitempty,min=1,max=1000" form:"height"`

	Published *bool `validate:"omitempty" form:"-"`
}

type ProjectAssign struct {
	View           *bool `validate:"omitempty" json:"view"`
	Edit           *bool `validate:"omitempty" json:"edit"`
	Delete         *bool `validate:"omitempty" json:"delete"`
	Publish        *bool `validate:"omitempty" json:"publish"`
	Share          *bool `validate:"omitempty" json:"share"`
	ManageUsers    *bool `validate:"omitempty" json:"manage_users"`
	ManageMetadata *bool `validate:"omitempty" json:"manage_metadata"`
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
	FirstName      string `json:"firstname"`
	LastName       string `json:"lastname"`
	Username       string `json:"username"`
	ProfilePicture string `json:"pfp"`
	Verified       bool   `json:"verified"`
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

	Width  int `json:"width"`
	Height int `json:"height"`

	BannerURL string `json:"banner_url"`

	IsPublic bool  `json:"is_public"`
	Fork     *uint `json:"fork"`

	OwnerID             uint   `json:"owner_id"`
	OwnerFirstName      string `json:"owner_firstname"`
	OwnerLastName       string `json:"owner_lastname"`
	OwnerUsername       string `json:"owner_username"`
	OwnerProfilePicture string `json:"owner_pfp"`
	OwnerVerified       bool   `json:"owner_verified"`

	AllowedUsers utils.JSON `gorm:"type:jsonb" json:"allowed_users"`

	IsFavorited    bool `json:"is_favorited"`
	TotalFavorites int  `json:"total_favorites"`
	TotalClones    int  `json:"total_clones"`
}

type ProjectInfo struct {
	ID        uint           `json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at"`

	Name        string `json:"name"`
	Description string `json:"description"`

	Budget int `json:"budget"`

	Width  int `json:"width"`
	Height int `json:"height"`

	BannerURL string `json:"banner_url"`

	IsPublic bool  `json:"is_public"`
	Fork     *uint `json:"fork"`

	OwnerID             uint   `json:"owner_id"`
	OwnerFirstName      string `json:"owner_firstname"`
	OwnerLastName       string `json:"owner_lastname"`
	OwnerUsername       string `json:"owner_username"`
	OwnerProfilePicture string `json:"owner_pfp"`
	OwnerVerified       bool   `json:"owner_verified"`

	AllowedUsers []ProjectUserPermissions `json:"allowed_users"`

	IsFavorited    bool `json:"is_favorited"`
	TotalFavorites int  `json:"total_favorites"`
	TotalClones    int  `json:"total_clones"`
}

func (a Allow) IsEmpty() bool {
	return !a.View && !a.Edit && !a.Delete && !a.Publish && !a.Share && !a.Manage.Users && !a.Manage.Metadata
}

func (a ProjectAssign) IsEmpty() bool {
	return ((a.View != nil && !*a.View) || a.View == nil) &&
		((a.Edit != nil && !*a.Edit) || a.Edit == nil) &&
		((a.Delete != nil && !*a.Delete) || a.Delete == nil) &&
		((a.Publish != nil && !*a.Publish) || a.Publish == nil) &&
		((a.Share != nil && !*a.Share) || a.Share == nil) &&
		((a.ManageUsers != nil && !*a.ManageUsers) || a.ManageUsers == nil) &&
		((a.ManageMetadata != nil && !*a.ManageMetadata) || a.ManageMetadata == nil)
}
