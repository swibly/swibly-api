package v1

import (
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/swibly/swibly-api/config"
	"github.com/swibly/swibly-api/internal/model/dto"
	"github.com/swibly/swibly-api/internal/service"
	"github.com/swibly/swibly-api/pkg/aws"
	"github.com/swibly/swibly-api/pkg/middleware"
	"github.com/swibly/swibly-api/pkg/utils"
	"github.com/swibly/swibly-api/translations"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func newAuthRoutes(handler *gin.RouterGroup) {
	h := handler.Group("/auth", middleware.APIKeyHasEnabledAuth)
	{
		h.GET("/validate", middleware.Auth, GetUserByBearerHandler)

		h.POST("/register", RegisterHandler)
		h.POST("/login", LoginHandler)

		h.PATCH("/update", middleware.APIKeyHasEnabledUserActions, middleware.Auth, UpdateUserHandler)
		h.PATCH("/image", middleware.APIKeyHasEnabledUserActions, middleware.Auth, UploadUserImage)

		h.DELETE("/delete", middleware.APIKeyHasEnabledUserActions, middleware.Auth, DeleteUserHandler)
		h.DELETE("/image", middleware.APIKeyHasEnabledUserActions, middleware.Auth, RemoveUserImage)
	}

	password := h.Group("/password")
	{
		password.POST("/reset", RequestPasswordResetHandler)
		password.POST("/reset/:key", PasswordResetHandler)

		password.OPTIONS("/reset/:key", ValidatePasswordResetHandler)
	}
}

func GetUserByBearerHandler(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, ctx.Keys["auth_user"])
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

	if body.Username == "" && body.Email == "" {
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

	issuer := ctx.Keys["auth_user"].(*dto.UserProfile)

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

	if body.Username != nil && *body.Username != "" && *body.Username != issuer.Username {
		if profile, err := service.User.GetByUsername(*body.Username); profile != nil && err == nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": dict.AuthDuplicatedUser})
			return
		}
	}

	if body.Email != nil && *body.Email != "" && *body.Email != issuer.Email {
		if profile, err := service.User.GetByEmail(*body.Email); profile != nil && err == nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": dict.AuthDuplicatedUser})
			return
		}
	}

	if body.Password != nil && *body.Password != "" {
		if hashedPassword, err := bcrypt.GenerateFromPassword([]byte(*body.Password), config.Security.BcryptCost); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": dict.InternalServerError})
		} else {
			*body.Password = string(hashedPassword)
		}
	}

	if err := service.User.Update(issuer.ID, &body); err != nil {
		log.Print(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": dict.InternalServerError})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": dict.AuthUserUpdated})
}

func UploadUserImage(ctx *gin.Context) {
	dict := translations.GetTranslation(ctx)

	issuer := ctx.Keys["auth_user"].(*dto.UserProfile)

	userProfilePicture := &dto.UserProfilePicture{}
	if err := ctx.Bind(userProfilePicture); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": dict.InvalidBody})
		return
	}

	if errs := utils.ValidateStruct(userProfilePicture); errs != nil {
		err := utils.ValidateErrorMessage(ctx, errs[0])

		log.Print(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": gin.H{err.Param: err.Message}})
		return
	}

	if err := service.User.SetProfilePicture(issuer, userProfilePicture.Image); err != nil {
		if errors.Is(err, aws.ErrUnsupportedFileType) {
			log.Print(err)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": dict.UnsupportedFileType})
			return
		}

		if errors.Is(err, aws.ErrUnableToDecode) {
			log.Print(err)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": dict.UnableToDecodeFile})
			return
		}

		if errors.Is(err, aws.ErrUnableToEncode) {
			log.Print(err)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": dict.UnableToDecodeFile})
			return
		}

		if errors.Is(err, aws.ErrFileTooLarge) {
			log.Print(err)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": dict.FileTooLarge})
			return
		}

		log.Print(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": dict.InternalServerError})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": dict.AuthUserUpdated})
}

func DeleteUserHandler(ctx *gin.Context) {
	dict := translations.GetTranslation(ctx)

	issuer := ctx.Keys["auth_user"].(*dto.UserProfile)

	if err := service.User.DeleteUser(issuer.ID); err != nil {
		log.Print(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": dict.InternalServerError})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": dict.AuthUserDeleted})
}

func RemoveUserImage(ctx *gin.Context) {
	dict := translations.GetTranslation(ctx)

	issuer := ctx.Keys["auth_user"].(*dto.UserProfile)

	if err := service.User.RemoveProfilePicture(issuer); err != nil {
		log.Print(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": dict.InternalServerError})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": dict.AuthUserUpdated})
}

func RequestPasswordResetHandler(ctx *gin.Context) {
	dict := translations.GetTranslation(ctx)

	var body dto.RequestPasswordReset

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

	if err := service.PasswordReset.Request(dict, body.Email); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Print(err)
			ctx.JSON(http.StatusForbidden, gin.H{"error": dict.UserNotFound})
			return
		}

		log.Print(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": dict.InternalServerError})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": dict.PasswordResetRequest})
}

func PasswordResetHandler(ctx *gin.Context) {
	dict := translations.GetTranslation(ctx)

	key := ctx.Param("key")

	var body dto.PasswordReset

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

	if err := service.PasswordReset.Reset(key, body.Password); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Print(err)
			ctx.JSON(http.StatusForbidden, gin.H{"error": dict.InvalidPasswordResetKey})
			return
		}

		log.Print(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": dict.InternalServerError})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": dict.PasswordResetSuccess})
}

func ValidatePasswordResetHandler(ctx *gin.Context) {
	dict := translations.GetTranslation(ctx)

	key := ctx.Param("key")

	passwordResetInfo, isValid, err := service.PasswordReset.IsKeyValid(key)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": dict.InternalServerError})
		return
	}

	if isValid {
		ctx.JSON(http.StatusAccepted, passwordResetInfo)
		return
	}

	ctx.JSON(http.StatusNotAcceptable, gin.H{"message": dict.InvalidPasswordResetKey})
}
