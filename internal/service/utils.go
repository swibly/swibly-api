package service

import (
	"github.com/gin-gonic/gin"
	"github.com/swibly/swibly-api/internal/model/dto"
)

func CreateNotification(ctx *gin.Context, createModel dto.CreateNotification, ids ...uint) error {
	notification, err := Notification.Create(createModel)
	if err != nil {
		return err
	}

	if err := Notification.SendToIDs(notification, ids); err != nil {
		return err
	}

	return nil
}
