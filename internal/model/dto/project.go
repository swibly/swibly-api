package dto

type ProjectCreation struct {
	Name        string `validate:"required,min=3,max=32"    json:"name"`
	Description string `validate:"omitempty,min=3,max=5000" json:"description"`

	Content any `validate:"omitempty" json:"content"`
	Budget  int `validate:"omitempty" json:"budget"`

	OwnerID uint `json:"-"` // Set using JWT

	Public bool `json:"-"` // Set in an URL query param
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
	UserInfoLite
	Allow Allow `json:"allow"`
}

type ProjectInfo struct {
	Name        string `json:"name"`
	Description string `json:"description"`

	Content any `json:"content"`
	Budget  int `json:"budget"`

	Public bool `json:"public"`

	Owner        UserInfoLite             `json:"owner"`
	AllowedUsers []ProjectUserPermissions `json:"allowed_users"`
}

func (a Allow) IsEmpty() bool {
	return !a.View && !a.Edit && !a.Delete && !a.Publish && !a.Share && !a.Manage.Users && !a.Manage.Metadata
}
