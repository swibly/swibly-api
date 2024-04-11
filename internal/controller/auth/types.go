package auth

import (
	"errors"

	"github.com/devkcud/arkhon-foundation/arkhon-api/internal/model"
	"github.com/devkcud/arkhon-foundation/arkhon-api/internal/utils"
	"golang.org/x/crypto/bcrypt"
)

var ErrUserAlreadyDefined error = errors.New("user already exists with that name or email")
var ErrInvalidUsernameOrEmail error = errors.New("not a valid username or email")

type UserBodyRegister struct {
	Fullname string `json:"fullname"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserBodyLogin struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func register(fullname, username, email, password string) (uint, error) {
	var exists bool
	if err := utils.DB.Model(&model.User{}).Select("count(*) > 0").Where("username = ? OR email = ?", username, email).Find(&exists).Error; err != nil || exists {
		return 0, err
	}

	if exists {
		return 0, ErrUserAlreadyDefined
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 10)

	if err != nil {
		return 0, nil
	}

	user := model.User{
		Fullname: fullname,
		Username: username,
		Email:    email,
		Password: string(hashedPassword),
	}

	if err := utils.DB.Create(&user).Error; err != nil {
		return 0, err
	}

	return user.ID, nil
}

func login(username, email, password string) (uint, error) {
	if username == "" && email == "" {
		return 0, ErrInvalidUsernameOrEmail
	}

	var user model.User

	if err := utils.DB.Model(&model.User{}).Where("username = ? OR email = ?", username, email).First(&user).Error; err != nil {
		return 0, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return 0, err
	}

	return user.ID, nil
}
