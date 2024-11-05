package repository

import (
	"errors"

	"github.com/swibly/swibly-api/internal/model"
	"github.com/swibly/swibly-api/internal/model/dto"
	"github.com/swibly/swibly-api/pkg/db"
	"github.com/swibly/swibly-api/pkg/pagination"
	"gorm.io/gorm"
)

type notificationRepository struct {
	db *gorm.DB
}

type NotificationRepository interface {
	Create(createModel dto.CreateNotification) (uint, error)

	GetForUser(userID uint, onlyUnread bool, page, perPage int) (*dto.Pagination[dto.NotificationInfo], error)
	GetUnreadCount(userID uint) (int64, error)

	SendToAll(notificationID uint) error
	SendToIDs(notificationID uint, usersID []uint) error

	UnsendToAll(notificationID uint) error
	UnsendToIDs(notificationID uint, usersID []uint) error

	MarkAsRead(issuer dto.UserProfile, notificationID uint) error
	MarkAsUnread(issuer dto.UserProfile, notificationID uint) error
}

var (
	ErrNotificationAlreadyRead = errors.New("notification already marked as read")
	ErrNotificationNotRead     = errors.New("notification is not marked as read")
	ErrNotificationNotAssigned = errors.New("notification not assigned to user")
)

func NewNotificationRepository() NotificationRepository {
	return &notificationRepository{db: db.Postgres}
}

func (nr *notificationRepository) baseNotificationQuery(userID uint) *gorm.DB {
	return nr.db.Table("notification_users AS nu").
		Select(`
			n.id AS id,
			n.created_at AS created_at,
			n.updated_at AS updated_at,
			n.title AS title,
			n.message AS message,
			n.type AS type,
			n.redirect AS redirect,
			nur.created_at AS read_at,
			CASE WHEN nur.created_at IS NOT NULL THEN true ELSE false END AS is_read
		`).
		Joins("JOIN notifications AS n ON n.id = nu.notification_id").
		Joins("LEFT JOIN notification_user_reads AS nur ON n.id = nur.notification_id AND nur.user_id = ?", userID).
		Where("nu.user_id = ?", userID)
}

func (nr *notificationRepository) Create(createModel dto.CreateNotification) (uint, error) {
	notification := &model.Notification{
		Title:    createModel.Title,
		Message:  createModel.Message,
		Type:     createModel.Type,
		Redirect: createModel.Redirect,
	}
	if err := nr.db.Create(notification).Error; err != nil {
		return 0, err
	}

	return notification.ID, nil
}

func (nr *notificationRepository) GetForUser(userID uint, onlyUnread bool, page, perPage int) (*dto.Pagination[dto.NotificationInfo], error) {
	query := nr.baseNotificationQuery(userID).Order("n.created_at DESC")

	if onlyUnread {
		query = query.Where("is_read IS false")
	}

	paginationResult, err := pagination.Generate[dto.NotificationInfo](query, page, perPage)
	if err != nil {
		return nil, err
	}

	return paginationResult, nil
}

func (nr *notificationRepository) GetUnreadCount(userID uint) (int64, error) {
	count := int64(0)
	if err := nr.baseNotificationQuery(userID).Where("is_read IS false").Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

func (nr *notificationRepository) SendToAll(notificationID uint) error {
	tx := nr.db.Begin()

	var userIDs []uint
	if err := tx.Model(&model.User{}).Select("id").Find(&userIDs).Error; err != nil {
		tx.Rollback()
		return err
	}

	notifications := make([]model.NotificationUser, len(userIDs))
	for i, userID := range userIDs {
		notifications[i] = model.NotificationUser{
			NotificationID: notificationID,
			UserID:         userID,
		}
	}

	if err := tx.Create(&notifications).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (nr *notificationRepository) SendToIDs(notificationID uint, userIDs []uint) error {
	tx := nr.db.Begin()

	notifications := make([]model.NotificationUser, len(userIDs))
	for i, userID := range userIDs {
		notifications[i] = model.NotificationUser{
			NotificationID: notificationID,
			UserID:         userID,
		}
	}

	if err := tx.Create(&notifications).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (nr *notificationRepository) UnsendToAll(notificationID uint) error {
	if err := nr.db.Where("notification_id = ?", notificationID).Delete(&model.NotificationUser{}).Error; err != nil {
		return err
	}

	if err := nr.db.Where("notification_id = ?", notificationID).Delete(&model.NotificationUserRead{}).Error; err != nil {
		return err
	}

	return nil
}

func (nr *notificationRepository) UnsendToIDs(notificationID uint, userIDs []uint) error {
	tx := nr.db.Begin()

	if len(userIDs) == 0 {
		return nil
	}

	if err := tx.Where("notification_id = ? AND user_id IN ?", notificationID, userIDs).Delete(&model.NotificationUser{}).Error; err != nil {
		return err
	}

	if err := tx.Where("notification_id = ? AND user_id IN ?", notificationID, userIDs).Delete(&model.NotificationUserRead{}).Error; err != nil {
		return err
	}

	return nil
}

func (nr *notificationRepository) MarkAsRead(issuer dto.UserProfile, notificationID uint) error {
	tx := nr.db.Begin()

	var notificationUser model.NotificationUser
	if err := tx.Where("notification_id = ? AND user_id = ?", notificationID, issuer.ID).First(&notificationUser).Error; err == gorm.ErrRecordNotFound {
		tx.Rollback()
		return ErrUserNotAssigned
	} else if err != nil {
		tx.Rollback()
		return err
	}

	var existingRead model.NotificationUserRead
	if err := tx.Where("notification_id = ? AND user_id = ?", notificationID, issuer.ID).First(&existingRead).Error; err == nil {
		tx.Rollback()
		return ErrNotificationAlreadyRead
	} else if err != gorm.ErrRecordNotFound {
		tx.Rollback()
		return err
	}

	if err := tx.Create(&model.NotificationUserRead{NotificationID: notificationID, UserID: issuer.ID}).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (nr *notificationRepository) MarkAsUnread(issuer dto.UserProfile, notificationID uint) error {
	tx := nr.db.Begin()

	var notificationUser model.NotificationUser
	if err := tx.Where("notification_id = ? AND user_id = ?", notificationID, issuer.ID).First(&notificationUser).Error; err == gorm.ErrRecordNotFound {
		tx.Rollback()
		return ErrUserNotAssigned
	} else if err != nil {
		tx.Rollback()
		return err
	}

	var existingRead model.NotificationUserRead
	if err := tx.Where("notification_id = ? AND user_id = ?", notificationID, issuer.ID).First(&existingRead).Error; err == gorm.ErrRecordNotFound {
		tx.Rollback()
		return ErrNotificationNotRead
	} else if err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Where("notification_id = ? AND user_id = ?", notificationID, issuer.ID).Delete(&model.NotificationUserRead{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}
