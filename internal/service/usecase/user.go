package usecase

import (
	"github.com/devkcud/arkhon-foundation/arkhon-api/internal/model"
	"github.com/devkcud/arkhon-foundation/arkhon-api/internal/service/repository"
)

type UserUseCase interface {
	CreateUser(firstname, lastname, username, email, password string) error
	DeleteUser(id uint) error

	GetByID(id uint) (*model.User, error)
	GetByUsername(username string) (*model.User, error)
	GetByEmail(email string) (*model.User, error)
	GetByUsernameOrEmail(username, email string) (*model.User, error)
}

type userUseCaso struct {
	ur repository.UserRepository
}

var UserUseCaseInstance = newUserUseCase()

func newUserUseCase() UserUseCase {
	return userUseCaso{ur: repository.UserRepositoryInstance}
}

func (uuc userUseCaso) CreateUser(firstname, lastname, username, email, password string) error {
	return nil
}

func (uuc userUseCaso) DeleteUser(id uint) error {
	return nil
}

func (uuc userUseCaso) GetByID(id uint) (*model.User, error) {
	return nil, nil
}

func (uuc userUseCaso) GetByUsername(username string) (*model.User, error) {
	return nil, nil
}

func (uuc userUseCaso) GetByEmail(email string) (*model.User, error) {
	return nil, nil
}

func (uuc userUseCaso) GetByUsernameOrEmail(username, email string) (*model.User, error) {
	return nil, nil
}
