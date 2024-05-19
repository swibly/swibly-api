package v1

import (
	"errors"
	"log"
	"net/http"

	"github.com/devkcud/arkhon-foundation/arkhon-api/internal/model/dto"
	"github.com/devkcud/arkhon-foundation/arkhon-api/internal/service/usecase"
	"github.com/devkcud/arkhon-foundation/arkhon-api/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgconn"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func newAuthRoutes(handler *gin.RouterGroup) {
	h := handler.Group("/auth")
	{
		h.POST("/register", RegisterHandler)
		h.POST("/login", LoginHandler)
	}
}

func RegisterHandler(ctx *gin.Context) {
	var body dto.UserRegister

	if err := ctx.BindJSON(&body); err != nil {
		log.Print(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Bad body format"})
		return
	}

	user, err := usecase.UserInstance.CreateUser(body.FirstName, body.LastName, body.Username, body.Email, body.Password)

	if err == nil {
		if token, err := utils.GenerateJWT(user.ID); err != nil {
			log.Print(err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error. Please, try again later."})
		} else {
			ctx.JSON(http.StatusOK, gin.H{"message": "User created", "token": token})
		}

		return
	}

	// Just print every time there is an error, no need to check what is the "context"
	log.Print(err)

	if validationErr, ok := err.(utils.ParamError); ok {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": gin.H{validationErr.Param: validationErr.Message}})
		return
	}

	var pgErr *pgconn.PgError
	// 23505 => duplicated key value violates unique constraint
	if errors.Is(err, gorm.ErrDuplicatedKey) || (errors.As(err, &pgErr) && pgErr.Code == "23505") {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "An user with that username or email already exists."})
		return
	}

	ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error. Please, try again later."})
}

func LoginHandler(ctx *gin.Context) {
	var body dto.UserLogin

	if err := ctx.BindJSON(&body); err != nil {
		log.Print(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Bad body format"})
		return
	}

	if errs := utils.ValidateStruct(&body); errs != nil {
		err := utils.ValidateErrorMessage(errs[0])

		log.Print(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": gin.H{err.Param: err.Message}})
		return
	}

	user, err := usecase.UserInstance.GetByUsernameOrEmail(body.Username, body.Email)

	if err != nil {
		log.Print(err)

		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "No user found with that username or email"})
			return
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error. Please, try again later."})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password)); err != nil {
		log.Print(err)

		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Password mismatch"})
			return
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error. Please, try again later."})
		return
	}

	if token, err := utils.GenerateJWT(user.ID); err != nil {
		log.Print(err)

		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error. Please, try again later."})
	} else {
		ctx.JSON(http.StatusOK, gin.H{"message": "User logged in", "token": token})
	}
}
