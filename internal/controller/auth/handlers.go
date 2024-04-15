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

	if err := ctx.BindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Bad body formatting"})
		return
	}

	var exists bool

	if err := utils.DB.Model(&model.User{}).Select("count(*) > 0").Where("username = ? OR email = ?", body.Username, body.Email).Find(&exists).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Something wrong happened in our servers. Try again later."})
		return
	}

	if exists {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "User already exists"})
		return
	}

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

	// NOTE: The hashing only runs after the validation for performance sake
	// and then we set the user.Password with the corrected password

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

	if err := utils.DB.Create(&user).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Something wrong happened in our servers. Try again later."})
		log.Println(err)
		return
	}

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

	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Bad JSON format."})
		return
	}

	if body.Username == "" && body.Email == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Either username or email must be provided."})
		return
	}

	var user model.User

	if err := utils.DB.Where("username = ? OR email = ?", body.Username, body.Email).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid credentials. Please check your email/username and password."})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Something wrong happened in our servers. Try again later."})
		}

		log.Println(err)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password)); err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid credentials. Please check your email/username and password."})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Something wrong happened in our servers. Try again later."})
		}

		log.Println(err)
		return
	}

	token, err := utils.GenerateJWT(user.ID)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Something wrong happened in our servers. Try again later."})
		log.Println(err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"token": token})
}
