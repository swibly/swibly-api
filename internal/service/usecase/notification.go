package usecase

import (
	"github.com/swibly/swibly-api/internal/model/dto"
	"github.com/swibly/swibly-api/internal/service/repository"
)

type NotificationUseCase struct {
	nr repository.NotificationRepository
}

func NewNotificationUseCase() NotificationUseCase {
	return NotificationUseCase{nr: repository.NewNotificationRepository()}
}

func (nuc *NotificationUseCase) Create(createModel dto.CreateNotification) (uint, error) {
	return nuc.nr.Create(createModel)
}

func (nuc *NotificationUseCase) GetForUser(userID uint, onlyUnread bool, page, perPage int) (*dto.Pagination[dto.NotificationInfo], error) {
	return nuc.nr.GetForUser(userID, onlyUnread, page, perPage)
}

func (nuc *NotificationUseCase) SendToAll(notificationID uint) error {
	return nuc.nr.SendToAll(notificationID)
}

func (nuc *NotificationUseCase) SendToIDs(notificationID uint, usersID []uint) error {
	return nuc.nr.SendToIDs(notificationID, usersID)
}

func (nuc *NotificationUseCase) UnsendToAll(notificationID uint) error {
	return nuc.nr.UnsendToAll(notificationID)
}

func (nuc *NotificationUseCase) UnsendToIDs(notificationID uint, usersID []uint) error {
	return nuc.nr.UnsendToIDs(notificationID, usersID)
}

func (nuc *NotificationUseCase) MarkAsRead(issuer dto.UserProfile, notificationID uint) error {
	return nuc.nr.MarkAsRead(issuer, notificationID)
}

func (nuc *NotificationUseCase) MarkAsUnread(issuer dto.UserProfile, notificationID uint) error {
	return nuc.nr.MarkAsUnread(issuer, notificationID)
}
