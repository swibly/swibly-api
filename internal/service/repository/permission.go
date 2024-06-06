package repository

import (
	"github.com/devkcud/arkhon-foundation/arkhon-api/internal/model"
	"github.com/devkcud/arkhon-foundation/arkhon-api/pkg/db"
	"gorm.io/gorm"
)

type permissionRepository struct {
	db *gorm.DB
}

type PermissionRepository interface {
	GetPermissions(userID uint) ([]*model.Permission, error)
}

func NewPermissionRepository() PermissionRepository {
	return permissionRepository{db: db.Postgres}
}

func (pr permissionRepository) GetPermissions(userID uint) ([]*model.Permission, error) {
	var permissions []*model.Permission

	err := pr.db.Table("users").
		Select("permissions.id, permissions.name").
		Joins("JOIN user_permissions ON user_permissions.user_id = users.id").
		Joins("JOIN permissions ON permissions.id = user_permissions.permission_id").
		Where("users.id = ?", userID).
		Scan(&permissions).Error

	if err != nil {
		return nil, err
	}

	return permissions, nil
}
