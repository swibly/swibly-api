package usecase

import (
	"github.com/devkcud/arkhon-foundation/arkhon-api/internal/model"
	"github.com/devkcud/arkhon-foundation/arkhon-api/internal/service/repository"
	"github.com/devkcud/arkhon-foundation/arkhon-api/pkg/utils"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserUseCase struct {
	ur repository.Repository[model.User]
}

func NewUserUseCase() UserUseCase {
	return UserUseCase{ur: repository.NewUserRepository()}
}

func (uuc UserUseCase) CreateUser(firstname, lastname, username, email, password string) error {
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
		newUser.Password = string(hashedPassword) // Set the hash
	}

	return uuc.ur.Store(&newUser)
}

func (uuc UserUseCase) DeleteUser(id uint) error {
	if _, err := uuc.GetByID(id); err != nil {
		return gorm.ErrRecordNotFound
	}

	return uuc.ur.Delete(id)
}

func (uuc UserUseCase) GetByID(id uint) (*model.User, error) {
	return uuc.ur.Find(&model.User{ID: id})
}

func (uuc UserUseCase) GetByUsername(username string) (*model.User, error) {
	return uuc.ur.Find(&model.User{Username: username})
}

func (uuc UserUseCase) GetByEmail(email string) (*model.User, error) {
	return uuc.ur.Find(&model.User{Email: email})
}

func (uuc UserUseCase) GetByUsernameOrEmail(username, email string) (*model.User, error) {
	return uuc.ur.Find(&model.User{Username: username, Email: email})
}
