package usecase

import (
	"github.com/devkcud/arkhon-foundation/arkhon-api/internal/model"
	"github.com/devkcud/arkhon-foundation/arkhon-api/internal/service/repository"
	"github.com/devkcud/arkhon-foundation/arkhon-api/pkg/utils"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type userUseCase struct {
	ur repository.Repository[model.User]
}

func NewUserUseCase() userUseCase {
	return userUseCase{ur: repository.NewUserRepository()}
}

func (uuc userUseCase) CreateUser(firstname, lastname, username, email, password string) error {
	newUser := model.User{
		FirstName: firstname,
		LastName:  lastname,
		Username:  username,
		Email:     email,
		Password:  password, // Hashing later
	}

	if errs := utils.ValidateStruct(&newUser); errs != nil {
		return utils.ValidateErrorMessage(errs[0])
	}

	if _, err := uuc.GetByUsernameOrEmail(username, email); err == nil {
		return gorm.ErrDuplicatedKey
	}

	if hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), 10); err != nil {
		return err
	} else {
		newUser.Password = string(hashedPassword)
	}

	return uuc.ur.Store(&newUser)
}

func (uuc userUseCase) DeleteUser(id uint) error {
	return nil
}

func (uuc userUseCase) GetByID(id uint) (*model.User, error) {
	return uuc.ur.Find(&model.User{ID: id})
}

func (uuc userUseCase) GetByUsername(username string) (*model.User, error) {
	return uuc.ur.Find(&model.User{Username: username})
}

func (uuc userUseCase) GetByEmail(email string) (*model.User, error) {
	return uuc.ur.Find(&model.User{Email: email})
}

func (uuc userUseCase) GetByUsernameOrEmail(username, email string) (*model.User, error) {
	return uuc.ur.Find(&model.User{Username: username, Email: email})
}
