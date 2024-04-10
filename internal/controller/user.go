package controller

import (
	"errors"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/devkcud/arkhon-foundation/arkhon-api/internal/loader"
	"github.com/devkcud/arkhon-foundation/arkhon-api/internal/model"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

func generateJWT(userID uint) (string, error) {
	// Generate JWT token
	// "sub" is the subject of the token (ID)
	// "exp" is the expiration date
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(time.Hour * 24 * 7).Unix(), // Will add 7 days
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))

	// Only related to internal error, not user input
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func register(fullname, username, email, password string) (uint, error) {
	var existingUser model.User
	if err := loader.DB.Model(&model.User{}).Where("username = ? OR email = ?", username, email).First(&existingUser).Error; err == nil {
		return 0, errors.New("username or email is already defined")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 10)

	if err != nil {
		return 0, err
	}

	user := model.User{
		Fullname: fullname,
		Username: username,
		Email:    email,
		Password: string(hashedPassword),
	}

	if err := loader.DB.Create(&user).Error; err != nil {
		return 0, err
	}

	return user.ID, nil
}

func login(username, email, password string) (uint, error) {
	if username == "" && email == "" {
		return 0, errors.New("not a valid username or email")
	}

	var user model.User

	if err := loader.DB.Model(&model.User{}).Where("username = ? OR email = ?", username, email).First(&user).Error; err != nil {
		return 0, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return 0, err
	}

	return user.ID, nil
}

func RegisterHandler(ctx *gin.Context) {
	var body struct {
		Fullname string `json:"fullname"`
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := ctx.BindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Bad body formatting"})
		return
	}

	userID, err := register(body.Fullname, body.Username, body.Email, body.Password)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		log.Println(err)
		return
	}

	token, err := generateJWT(userID)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		log.Println(err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"token": token})
}

func LoginHandler(ctx *gin.Context) {
	var body struct {
		Username string `json:"username" default:""`
		Email    string `json:"email" default:""`
		Password string `json:"password" default:""`
	}

	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Bad body formatting"})
		return
	}

	userID, err := login(body.Username, body.Email, body.Password)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		log.Println(err)
		return
	}

	token, err := generateJWT(userID)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		log.Println(err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"token": token})
}
