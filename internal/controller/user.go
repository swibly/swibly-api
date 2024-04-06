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
	"gorm.io/gorm"
)

func Register(ctx *gin.Context) {
	var body struct {
		Fullname string
		Username string
		Email    string
		Password string
	}

	if ctx.BindJSON(&body) != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Bad body formatting"})
		return
	}

	// Check if the user exists by getting the first entry in the database and checking if it's an error
	var existingUser model.User
	if err := loader.DB.Where("username = ? OR email = ?", body.Username, body.Email).First(&existingUser).Error; err == nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "User already exists"})
		return
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		log.Println("Failed to query database:", err)
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)

	// This is related to GenerateFromPassword (not user input) error
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		log.Println(err)
		return
	}

	user := model.User{Fullname: body.Fullname, Username: body.Username, Email: body.Email, Password: string(hash)}
	result := loader.DB.Create(&user)

	// Again, this is related only to the error of the db insertion, not the user input itself
	if result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		log.Println("Couldn't insert user in database", result.Error)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "User created"})
}

func Login(ctx *gin.Context) {
	// FIXME: Email is the only method available, should be Email and Username
	var body struct {
		Email    string
		Password string
	}

	if ctx.BindJSON(&body) != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Bad body formatting"})
		return
	}

	var user model.User
	// Check user two times, first by email, then by password
	if err := loader.DB.Where("email = ?", body.Email).First(&user).Error; err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		log.Println("Failed to find user by email:", err)
		return
	}

	// Compare hashed password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password)); err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "No user found with that email or password"})
		log.Println("Couldn't get by hash and password:", err)
		return
	}

	// Generate JWT token
	// "sub" is the subject of the token (user.ID)
	// "exp" is the expiration date
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 24 * 7).Unix(), // Will add 7 days
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))

	// Same thing, only related to internal error, not user input
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		log.Println("Couldn't generate token:", err)
		return
	}

	// TODO: Update the options in the cookie to match the actual path, domain and secure (all mandatory for security sake)
	ctx.SetSameSite(http.SameSiteLaxMode)
	ctx.SetCookie("Authorization", tokenString, 3600*24*7, "", "", false, true) // 3600*24*7 = 7 days

	ctx.JSON(http.StatusOK, gin.H{"token": tokenString})
}
