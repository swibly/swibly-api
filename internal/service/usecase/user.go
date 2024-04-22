package usecase

import (
	"github.com/devkcud/arkhon-foundation/arkhon-api/internal/model"
	"github.com/devkcud/arkhon-foundation/arkhon-api/internal/service/repository"
)

type userUseCase struct {
	ur repository.Repository[model.User]
}

func NewUserUseCase() userUseCase {
	return userUseCase{ur: repository.NewUserRepository()}
}

func (uuc userUseCase) CreateUser(firstname, lastname, username, email, password string) error {
	return nil
}

func (uuc userUseCase) DeleteUser(id uint) error {
	return nil
}

func (uuc userUseCase) GetByID(id uint) (*model.User, error) {
	return nil, nil
}

func (uuc userUseCase) GetByUsername(username string) (*model.User, error) {
	return nil, nil
}

func (uuc userUseCase) GetByEmail(email string) (*model.User, error) {
	return nil, nil
}

func (uuc userUseCase) GetByUsernameOrEmail(username, email string) (*model.User, error) {
	return nil, nil
}
