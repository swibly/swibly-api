package usecase

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"mime/multipart"

	"github.com/gin-gonic/gin"
	"github.com/swibly/swibly-api/config"
	"github.com/swibly/swibly-api/internal/model"
	"github.com/swibly/swibly-api/internal/model/dto"
	"github.com/swibly/swibly-api/internal/service/repository"
	"github.com/swibly/swibly-api/pkg/aws"
	"github.com/swibly/swibly-api/pkg/utils"
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

	hasher := sha256.Sum256([]byte(email))
	newUser.ProfilePicture = fmt.Sprintf("https://www.gravatar.com/avatar/%s?s=512&d=monsterid&r=g", hex.EncodeToString(hasher[:]))

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

	if err := uuc.ur.Create(&newUser); err != nil {
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
	return uuc.ur.Get(&model.User{ID: id})
}

func (uuc UserUseCase) GetByUsername(username string) (*dto.UserProfile, error) {
	return uuc.ur.Get(&model.User{Username: username})
}

func (uuc UserUseCase) GetByEmail(email string) (*dto.UserProfile, error) {
	return uuc.ur.Get(&model.User{Email: email})
}

func (uuc UserUseCase) GetByUsernameOrEmail(username, email string) (*dto.UserProfile, error) {
	return uuc.ur.Get(&model.User{Username: username, Email: email})
}

func (uuc UserUseCase) Search(issuerID uint, search *dto.SearchUser, page, perpage int) (*dto.Pagination[dto.UserProfile], error) {
	return uuc.ur.Search(issuerID, search, page, perpage)
}

func (uuc UserUseCase) UnsafeGetByID(id uint) (*model.User, error) {
	return uuc.ur.UnsafeGet(&model.User{ID: id})
}

func (uuc UserUseCase) UnsafeGetByUsername(username string) (*model.User, error) {
	return uuc.ur.UnsafeGet(&model.User{Username: username})
}

func (uuc UserUseCase) UnsafeGetByEmail(email string) (*model.User, error) {
	return uuc.ur.UnsafeGet(&model.User{Email: email})
}

func (uuc UserUseCase) UnsafeGetByUsernameOrEmail(username, email string) (*model.User, error) {
	return uuc.ur.UnsafeGet(&model.User{Username: username, Email: email})
}

func (uuc UserUseCase) Update(id uint, newModel *dto.UserUpdate) error {
	return uuc.ur.Update(id, newModel)
}

func (uuc UserUseCase) SetProfilePicture(user *dto.UserProfile, file *multipart.FileHeader) error {
	if err := aws.DeleteUserImage(user.ProfilePicture); err != nil {
		return err
	}

	url, err := aws.UploadUserImage(user.ID, file)
	if err != nil {
		return err
	}

	return uuc.Update(user.ID, &dto.UserUpdate{
		ProfilePicture: &url,
	})
}

func (uuc UserUseCase) RemoveProfilePicture(user *dto.UserProfile) error {
	hasher := sha256.Sum256([]byte(user.Email))

	if err := aws.DeleteUserImage(user.ProfilePicture); err != nil {
		return err
	}

	return uuc.Update(user.ID, &dto.UserUpdate{
		ProfilePicture: utils.ToPtr(fmt.Sprintf("https://www.gravatar.com/avatar/%s?s=512&d=monsterid&r=g", hex.EncodeToString(hasher[:]))),
	})
}
