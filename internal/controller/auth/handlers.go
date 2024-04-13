package auth

import (
	"log"
	"net/http"

	"github.com/devkcud/arkhon-foundation/arkhon-api/internal/model"
	"github.com/devkcud/arkhon-foundation/arkhon-api/internal/utils"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func RegisterHandler(ctx *gin.Context) {
	var body UserBodyRegister

	// Attempt to bind the request body to the UserBodyRegister struct
	if err := ctx.BindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Bad body formatting"})
		return
	}

	var exists bool
	// Check if user already exists by username or email
	if err := utils.DB.Model(&model.User{}).Select("count(*) > 0").Where("username = ? OR email = ?", body.Username, body.Email).Find(&exists).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Something wrong happened in our servers. Try again later."})
		return
	}

	// Check if the user already exists
	if exists {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "User already exists"})
		return
	}

	// Create a new User model with the provided and hashed information
	user := model.User{
		Fullname: body.Fullname,
		Username: body.Username,
		Email:    body.Email,
		Password: body.Password,
	}

	buildErrs, err := user.Validate()

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Something wrong happened in our servers. Try again later."})
		log.Println(err)
		return
	}

	if buildErrs != nil {
		ctx.JSON(http.StatusBadRequest, buildErrs)
		return
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)

	if err != nil {
		if err == bcrypt.ErrPasswordTooLong {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "The password is too long"})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Something wrong happened in our servers. Try again later."})
		}

		log.Println(err)
		return
	}

	user.Password = string(hashedPassword)

	// Attempt to create the user record in the database
	if err := utils.DB.Create(&user).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Something wrong happened in our servers. Try again later."})
		log.Println(err)
		return
	}

	// Attempt to generate a JWT token for the new user
	token, err := utils.GenerateJWT(user.ID)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		log.Println(err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"token": token})
}

func LoginHandler(ctx *gin.Context) {
	var body UserBodyLogin

	// Attempt to bind the incoming JSON to the struct
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Bad JSON format."})
		return
	}

	// Ensure either username or email is provided for authentication
	if body.Username == "" && body.Email == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Either username or email must be provided."})
		return
	}

	var user model.User

	// Check if a user exists with the given username or email
	if err := utils.DB.Where("username = ? OR email = ?", body.Username, body.Email).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid credentials. Please check your email/username and password."})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Something wrong happened in our servers. Try again later."})
		}

		log.Println(err)
		return
	}

	// Authenticate the user by comparing the hashed password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password)); err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid credentials. Please check your email/username and password."})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Something wrong happened in our servers. Try again later."})
		}

		log.Println(err)
		return
	}

	// Generate JWT token for authenticated user
	token, err := utils.GenerateJWT(user.ID)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Something wrong happened in our servers. Try again later."})
		log.Println(err)
		return
	}

	// Return the generated token to the user
	ctx.JSON(http.StatusOK, gin.H{"token": token})
}
