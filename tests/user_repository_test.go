package tests

import (
	"testing"

	"github.com/devkcud/arkhon-foundation/arkhon-api/internal/model"
	"github.com/devkcud/arkhon-foundation/arkhon-api/internal/service/repository"
	"go.uber.org/mock/gomock"
)

// NOTE: Those tests are just for the repository mocks, not actually storing anything or testing real things.

var userMock = &model.User{
	FirstName: "Test",
	LastName:  "Subject",
	Username:  "testsubject",
	Email:     "test.subject@example.com",
	Password:  "t&ST!ngF3@tur3",
}

func TestUserRepository_Store(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repository.NewMockUserRepository(ctrl)

	mockRepo.EXPECT().Store(userMock).Return(nil)

	if err := mockRepo.Store(userMock); err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestUserRepository_Update(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repository.NewMockUserRepository(ctrl)
	userToUpdate := &model.User{
		FirstName: "Updated",
		LastName:  "Name",
		Username:  "updatedname",
		Email:     "updated.name@example.com",
		Password:  "updatedPassword",
	}

	id := uint(1)

	mockRepo.EXPECT().Update(id, userToUpdate).Return(nil)

	err := mockRepo.Update(id, userToUpdate)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestUserRepository_Find(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repository.NewMockUserRepository(ctrl)
	searchModel := &model.User{
		FirstName: "Test",
		LastName:  "Subject",
	}

	mockRepo.EXPECT().Find(searchModel).Return(&model.User{
		FirstName: "Test",
		LastName:  "Subject",
		Username:  "testsubject",
		Email:     "test.subject@example.com",
		Password:  "t&STngF3@tur3",
	}, nil)

	foundUser, err := mockRepo.Find(searchModel)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if foundUser == nil || foundUser.FirstName != searchModel.FirstName ||
		foundUser.LastName != searchModel.LastName {
		t.Errorf("found user does not match expected user")
	}
}

func TestUserRepository_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repository.NewMockUserRepository(ctrl)

	id := uint(1)

	mockRepo.EXPECT().Delete(id).Return(nil)

	err := mockRepo.Delete(id)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}
