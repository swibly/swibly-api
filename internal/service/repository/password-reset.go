package repository

import (
	"time"

	"github.com/google/uuid"
	"github.com/swibly/swibly-api/config"
	"github.com/swibly/swibly-api/internal/model"
	"github.com/swibly/swibly-api/internal/model/dto"
	"github.com/swibly/swibly-api/pkg/db"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type passwordResetRepository struct {
	db *gorm.DB

	userRepo UserRepository
}

type PasswordResetRepository interface {
	Request(email string) (*model.PasswordResetKey, error)
	Reset(key, newPassword string) error
	IsKeyValid(key string) (*dto.PasswordResetInfo, bool, error)
}

func NewPasswordResetRepository() PasswordResetRepository {
	return &passwordResetRepository{db: db.Postgres, userRepo: NewUserRepository()}
}

func (prr *passwordResetRepository) Request(email string) (*model.PasswordResetKey, error) {
	user, err := prr.userRepo.Get(&model.User{Email: email})
	if err != nil {
		return nil, err
	}

	tx := prr.db.Begin()

	passwordReset := &model.PasswordResetKey{}
	err = tx.Where("user_id = ?", user.ID).First(passwordReset).Error

	if err == gorm.ErrRecordNotFound {
		passwordReset = &model.PasswordResetKey{
			UserID:    user.ID,
			Key:       uuid.New(),
			ExpiresAt: time.Now().Add(24 * time.Hour),
		}
		if err := tx.Create(passwordReset).Error; err != nil {
			tx.Rollback()
			return nil, err
		}
	} else if err == nil {
		passwordReset.Key = uuid.New()
		passwordReset.ExpiresAt = time.Now().Add(24 * time.Hour)

		if err := tx.Model(passwordReset).UpdateColumns(map[string]interface{}{
			"key":        passwordReset.Key,
			"expires_at": passwordReset.ExpiresAt,
		}).Error; err != nil {
			tx.Rollback()
			return nil, err
		}
	} else {
		tx.Rollback()
		return nil, err
	}

	return passwordReset, tx.Commit().Error
}

func (prr *passwordResetRepository) Reset(key, newPassword string) error {
	passwordReset := model.PasswordResetKey{}
	if err := prr.db.Model(&model.PasswordResetKey{}).Where("key = ?", key).Scan(&passwordReset).Error; err != nil {
		return err
	}

	if passwordReset.ExpiresAt.Before(time.Now()) {
		return gorm.ErrRecordNotFound
	}

	if _, err := prr.userRepo.UnsafeGet(&model.User{ID: passwordReset.UserID}); err != nil {
		return err
	}

	password := newPassword

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), config.Security.BcryptCost)
	if err != nil {
		return err
	}
	password = string(hashedPassword)

	tx := prr.db.Begin()

	if err := tx.Model(&model.User{}).Where("id = ?", passwordReset.UserID).Update("password", password).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Where("key = ?", key).Delete(&model.PasswordResetKey{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (prr *passwordResetRepository) IsKeyValid(key string) (*dto.PasswordResetInfo, bool, error) {
	passwordReset := model.PasswordResetKey{}

	if err := prr.db.Model(&model.PasswordResetKey{}).Where("key = ?", key).Scan(&passwordReset).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, false, nil
		}
		return nil, false, err
	}

	if passwordReset.ExpiresAt.Before(time.Now()) {
		return nil, false, nil
	}

	user := model.User{}
	if err := prr.db.Model(&model.User{}).Where("id = ?", passwordReset.UserID).Scan(&user).Error; err != nil {
		return nil, false, err
	}

	passwordResetInfo := &dto.PasswordResetInfo{
		FirstName:      user.FirstName,
		LastName:       user.LastName,
		Username:       user.Username,
		ProfilePicture: user.ProfilePicture,
		Lang:           string(user.Language),
		ExpiresAt:      passwordReset.ExpiresAt,
	}

	return passwordResetInfo, true, nil
}
