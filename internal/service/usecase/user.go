package usecase

import (
	"github.com/devkcud/arkhon-foundation/arkhon-api/config"
	"github.com/devkcud/arkhon-foundation/arkhon-api/internal/model"
	"github.com/devkcud/arkhon-foundation/arkhon-api/internal/model/dto"
	"github.com/devkcud/arkhon-foundation/arkhon-api/internal/service/repository"
	"github.com/devkcud/arkhon-foundation/arkhon-api/pkg/utils"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserUseCase struct {
	ur repository.UserRepository
}

func NewUserUseCase() UserUseCase {
	return UserUseCase{ur: repository.NewUserRepository()}
}

func (uuc UserUseCase) CreateUser(ctx *gin.Context, firstname, lastname, username, email, password string) (*model.User, error) {
	newUser := model.User{
		FirstName: firstname,
		LastName:  lastname,
		Username:  username,
		Email:     email,
		Password:  password, // Hashing later
	}

	if errs := utils.ValidateStruct(&newUser); errs != nil {
		return nil, utils.ValidateErrorMessage(ctx, errs[0])
	}

	if _, err := uuc.GetByUsernameOrEmail(username, email); err == nil {
		return nil, gorm.ErrDuplicatedKey
	}

	if hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), config.Security.BcryptCost); err != nil {
		return nil, err
	} else {
		newUser.Password = string(hashedPassword) // Set the hash
	}

	if err := uuc.ur.Store(&newUser); err != nil {
		return nil, err
	}

	return &newUser, nil
}

func (uuc UserUseCase) DeleteUser(id uint) error {
	if _, err := uuc.GetByID(id); err != nil {
		return gorm.ErrRecordNotFound
	}

	return uuc.ur.Delete(id)
}

func (uuc UserUseCase) GetByID(id uint) (*dto.UserProfile, error) {
	return uuc.ur.Find(&model.User{ID: id})
}

func (uuc UserUseCase) GetByUsername(username string) (*dto.UserProfile, error) {
	return uuc.ur.Find(&model.User{Username: username})
}

func (uuc UserUseCase) GetByEmail(email string) (*dto.UserProfile, error) {
	return uuc.ur.Find(&model.User{Email: email})
}

func (uuc UserUseCase) GetByUsernameOrEmail(username, email string) (*dto.UserProfile, error) {
	return uuc.ur.Find(&model.User{Username: username, Email: email})
}

func (uuc UserUseCase) GetBySimilarName(name string, page, perpage int) (*dto.Pagination[dto.UserProfile], error) {
	return uuc.ur.SearchLikeName(name, page, perpage)
}

func (uuc UserUseCase) UnsafeGetByID(id uint) (*model.User, error) {
	return uuc.ur.UnsafeFind(&model.User{ID: id})
}

func (uuc UserUseCase) UnsafeGetByUsername(username string) (*model.User, error) {
	return uuc.ur.UnsafeFind(&model.User{Username: username})
}

func (uuc UserUseCase) UnsafeGetByEmail(email string) (*model.User, error) {
	return uuc.ur.UnsafeFind(&model.User{Email: email})
}

func (uuc UserUseCase) UnsafeGetByUsernameOrEmail(username, email string) (*model.User, error) {
	return uuc.ur.UnsafeFind(&model.User{Username: username, Email: email})
}

func (uuc UserUseCase) Update(id uint, newModel *dto.UserUpdate) error {
	return uuc.ur.Update(id, newModel)
}
