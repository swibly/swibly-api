package v1

import (
	"errors"
	"log"
	"net/http"

	"github.com/devkcud/arkhon-foundation/arkhon-api/config"
	"github.com/devkcud/arkhon-foundation/arkhon-api/internal/model/dto"
	"github.com/devkcud/arkhon-foundation/arkhon-api/internal/service"
	"github.com/devkcud/arkhon-foundation/arkhon-api/pkg/middleware"
	"github.com/devkcud/arkhon-foundation/arkhon-api/pkg/utils"
	"github.com/devkcud/arkhon-foundation/arkhon-api/translations"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgconn"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func newAuthRoutes(handler *gin.RouterGroup) {
	h := handler.Group("/auth")
	h.Use(middleware.APIKeyHasEnabledAuth)
	{
		h.GET("/validate", middleware.AuthMiddleware) // Made so you can look up whenever you want to check if the token is valid
		h.POST("/register", RegisterHandler)
		h.POST("/login", LoginHandler)
		h.PATCH("/update", middleware.AuthMiddleware, UpdateUserHandler)
		h.DELETE("/delete", middleware.AuthMiddleware, DeleteUserHandler)
	}
}

func RegisterHandler(ctx *gin.Context) {
	dict := translations.GetTranslation(ctx)

	var body dto.UserRegister

	if err := ctx.BindJSON(&body); err != nil {
		log.Print(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": dict.InvalidBody})
		return
	}

	if errs := utils.ValidateStruct(&body); errs != nil {
		err := utils.ValidateErrorMessage(ctx, errs[0])

		log.Print(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": gin.H{err.Param: err.Message}})
		return
	}

	user, err := service.User.CreateUser(ctx, body.FirstName, body.LastName, body.Username, body.Email, body.Password)

	if err == nil {
		if token, err := utils.GenerateJWT(user.ID); err != nil {
			log.Print(err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": dict.InternalServerError})
		} else {
			ctx.JSON(http.StatusOK, gin.H{"token": token})
		}

		return
	}

	log.Print(err)

	if validationErr, ok := err.(utils.ParamError); ok {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": gin.H{validationErr.Param: validationErr.Message}})
		return
	}

	var pgErr *pgconn.PgError
	// 23505 => duplicated key value violates unique constraint
	if errors.Is(err, gorm.ErrDuplicatedKey) || (errors.As(err, &pgErr) && pgErr.Code == "23505") {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": dict.AuthDuplicatedUser})
		return
	}

	ctx.JSON(http.StatusInternalServerError, gin.H{"error": dict.InternalServerError})
}

func LoginHandler(ctx *gin.Context) {
	dict := translations.GetTranslation(ctx)

	var body dto.UserLogin

	if err := ctx.BindJSON(&body); err != nil {
		log.Print(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": dict.InvalidBody})
		return
	}

	if errs := utils.ValidateStruct(&body); errs != nil {
		err := utils.ValidateErrorMessage(ctx, errs[0])

		log.Print(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": gin.H{err.Param: err.Message}})
		return
	}

	user, err := service.User.UnsafeGetByUsernameOrEmail(body.Username, body.Email)

	if err != nil {
		log.Print(err)

		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": dict.AuthWrongCredentials})
			return
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{"error": dict.InternalServerError})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password)); err != nil {
		log.Print(err)

		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": dict.AuthWrongCredentials})
			return
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{"error": dict.InternalServerError})
		return
	}

	if token, err := utils.GenerateJWT(user.ID); err != nil {
		log.Print(err)

		ctx.JSON(http.StatusInternalServerError, gin.H{"error": dict.InternalServerError})
	} else {
		ctx.JSON(http.StatusOK, gin.H{"token": token})
	}
}

func UpdateUserHandler(ctx *gin.Context) {
	dict := translations.GetTranslation(ctx)

	issuer := ctx.Keys["auth_user"].(*dto.ProfileSearch)

	var body dto.UserUpdate

	if err := ctx.BindJSON(&body); err != nil {
		log.Print(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": dict.InvalidBody})
		return
	}

	if errs := utils.ValidateStruct(&body); errs != nil {
		err := utils.ValidateErrorMessage(ctx, errs[0])

		log.Print(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": gin.H{err.Param: err.Message}})
		return
	}

	if body.Username != "" {
		if profile, err := service.User.GetByUsername(body.Username); profile != nil && err == nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": dict.AuthDuplicatedUser})
			return
		}
	}

	if body.Email != "" {
		if profile, err := service.User.GetByEmail(body.Email); profile != nil && err == nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": dict.AuthDuplicatedUser})
			return
		}
	}

	if body.Password != "" {
		if hashedPassword, err := bcrypt.GenerateFromPassword([]byte(body.Password), config.Security.BcryptCost); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": dict.InternalServerError})
		} else {
			body.Password = string(hashedPassword)
		}
	}

	if err := service.User.Update(issuer.ID, &body); err != nil {
		log.Print(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": dict.InternalServerError})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": dict.AuthUserUpdated})
}

func DeleteUserHandler(ctx *gin.Context) {
	dict := translations.GetTranslation(ctx)

	issuer := ctx.Keys["auth_user"].(*dto.ProfileSearch)

	if err := service.User.DeleteUser(issuer.ID); err != nil {
		log.Print(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": dict.InternalServerError})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": dict.AuthUserDeleted})
}
