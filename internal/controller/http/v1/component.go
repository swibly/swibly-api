package v1

import (
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/devkcud/arkhon-foundation/arkhon-api/internal/model/dto"
	"github.com/devkcud/arkhon-foundation/arkhon-api/internal/service"
	"github.com/devkcud/arkhon-foundation/arkhon-api/internal/service/repository"
	"github.com/devkcud/arkhon-foundation/arkhon-api/pkg/middleware"
	"github.com/devkcud/arkhon-foundation/arkhon-api/pkg/utils"
	"github.com/devkcud/arkhon-foundation/arkhon-api/translations"
	"github.com/gin-gonic/gin"
)

func newComponentRoutes(handler *gin.RouterGroup) {
	h := handler.Group("/components", middleware.APIKeyHasEnabledProjects, middleware.Auth)
	{
		h.GET("", GetPublicComponentsHandler)
		h.GET("/trash", GetTrashComponentsHandler)

		h.POST("", CreateComponentHandler)

		h.DELETE("/trash", DeleteTrashComponentsHandler)

		byUser := h.Group("/user/:username", middleware.UserLookup)
		{
			byUser.GET("", GetComponentsByUserHandler)
			byUser.GET("/owned", GetOwnedComponentsByUserHandler)
		}
	}

	specific := h.Group("/:id", middleware.ComponentLookup)
	{
		specific.GET("", GetComponentHandler)

		specific.POST("/buy", BuyComponentHandler)
		specific.POST("/sell", SellComponentHandler)

		specific.PATCH("/update", middleware.ComponentOwnership, UpdateComponentHandler)
		specific.PATCH("/publish", middleware.ComponentOwnership, PublishComponentHandler)

		specific.DELETE("/unpublish", middleware.ComponentOwnership, UnpublishComponentHandler)

		trashActions := specific.Group("/trash", middleware.ComponentOwnership)
		{
			trashActions.PATCH("/restore", RestoreComponentHandler)

			trashActions.DELETE("", DeleteComponentHandler)
			trashActions.DELETE("/force", DeleteComponenttForceHandler)
		}
	}
}

func GetPublicComponentsHandler(ctx *gin.Context) {
	dict := translations.GetTranslation(ctx)

	issuerID := ctx.Keys["auth_user"].(*dto.UserProfile).ID

	page := 1
	perPage := 10
	freeOnly := false

	if i, e := strconv.Atoi(ctx.Query("page")); e == nil && ctx.Query("page") != "" {
		page = i
	}

	if i, e := strconv.Atoi(ctx.Query("perpage")); e == nil && ctx.Query("perpage") != "" {
		perPage = i
	}

	freeStatus := strings.ToLower(ctx.Query("free"))
	if freeStatus == "true" || freeStatus == "t" || freeStatus == "1" {
		freeOnly = true
	}

	components, err := service.Component.GetPublic(issuerID, page, perPage, freeOnly)
	if err != nil {
		log.Print(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": dict.InternalServerError})
		return
	}

	ctx.JSON(http.StatusOK, components)
}

func GetTrashComponentsHandler(ctx *gin.Context) {
	dict := translations.GetTranslation(ctx)

	issuerID := ctx.Keys["auth_user"].(*dto.UserProfile).ID

	page := 1
	perPage := 10

	if i, e := strconv.Atoi(ctx.Query("page")); e == nil && ctx.Query("page") != "" {
		page = i
	}

	if i, e := strconv.Atoi(ctx.Query("perpage")); e == nil && ctx.Query("perpage") != "" {
		perPage = i
	}

	components, err := service.Component.GetTrashed(issuerID, page, perPage)
	if err != nil {
		log.Print(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": dict.InternalServerError})
		return
	}

	ctx.JSON(http.StatusOK, components)
}

func CreateComponentHandler(ctx *gin.Context) {
	dict := translations.GetTranslation(ctx)

	issuer := ctx.Keys["auth_user"].(*dto.UserProfile)

	component := &dto.ComponentCreation{}
	if err := ctx.BindJSON(component); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if errs := utils.ValidateStruct(component); errs != nil {
		err := utils.ValidateErrorMessage(ctx, errs[0])

		log.Print(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": gin.H{err.Param: err.Message}})
		return
	}

	publicStatus := strings.ToLower(ctx.Query("public"))
	if publicStatus == "true" || publicStatus == "t" || publicStatus == "1" {
		component.Public = true
	}

	component.OwnerID = issuer.ID

	if err := service.Component.Create(component); err != nil {
		log.Print(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": dict.InternalServerError})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": dict.ComponentCreated})
}

func DeleteTrashComponentsHandler(ctx *gin.Context) {
	dict := translations.GetTranslation(ctx)

	service.Component.ClearTrash(ctx.Keys["auth_user"].(*dto.UserProfile).ID)

	ctx.JSON(http.StatusOK, gin.H{"message": dict.TrashCleared})
}

func GetComponentsByUserHandler(ctx *gin.Context) {
	dict := translations.GetTranslation(ctx)

	issuer := ctx.Keys["auth_user"].(*dto.UserProfile)
	user := ctx.Keys["user_lookup"].(*dto.UserProfile)

	page := 1
	perPage := 10

	if i, e := strconv.Atoi(ctx.Query("page")); e == nil && ctx.Query("page") != "" {
		page = i
	}

	if i, e := strconv.Atoi(ctx.Query("perpage")); e == nil && ctx.Query("perpage") != "" {
		perPage = i
	}

	components, err := service.Component.GetByOwnerID(issuer.ID, user.ID, issuer.ID != user.ID, page, perPage)
	if err != nil {
		log.Print(err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": dict.InternalServerError})
		return
	}

	ctx.JSON(http.StatusOK, components)
}

func GetOwnedComponentsByUserHandler(ctx *gin.Context) {
	dict := translations.GetTranslation(ctx)

	issuer := ctx.Keys["auth_user"].(*dto.UserProfile)
	user := ctx.Keys["user_lookup"].(*dto.UserProfile)

	page := 1
	perPage := 10

	if i, e := strconv.Atoi(ctx.Query("page")); e == nil && ctx.Query("page") != "" {
		page = i
	}

	if i, e := strconv.Atoi(ctx.Query("perpage")); e == nil && ctx.Query("perpage") != "" {
		perPage = i
	}

	components, err := service.Component.GetOwned(issuer.ID, user.ID, issuer.ID != user.ID, page, perPage)
	if err != nil {
		log.Print(err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": dict.InternalServerError})
		return
	}

	ctx.JSON(http.StatusOK, components)
}

func GetComponentHandler(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, ctx.Keys["component_lookup"].(*dto.ComponentInfo))
}

func UpdateComponentHandler(ctx *gin.Context) {
	dict := translations.GetTranslation(ctx)

	component := ctx.Keys["component_lookup"].(*dto.ComponentInfo)

	var body *dto.ComponentUpdate
	if err := ctx.BindJSON(&body); err != nil {
		log.Print(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": dict.InvalidBody})
		return
	}

	if errs := utils.ValidateStruct(body); errs != nil {
		err := utils.ValidateErrorMessage(ctx, errs[0])

		log.Print(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": gin.H{err.Param: err.Message}})
		return
	}

	if err := service.Component.Update(component.ID, body); err != nil {
		if errors.Is(err, repository.ErrComponentNotFound) {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": dict.ComponentAlreadyTrashed})
			return
		}

		log.Print(err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": dict.InternalServerError})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": dict.ComponentUpdated})
}

func BuyComponentHandler(ctx *gin.Context) {
	dict := translations.GetTranslation(ctx)

	issuerID := ctx.Keys["auth_user"].(*dto.UserProfile).ID
	componentID := ctx.Keys["component_lookup"].(*dto.ComponentInfo).ID

	if err := service.Component.Buy(issuerID, componentID); err != nil {
		if errors.Is(err, repository.ErrInsufficientArkhoins) {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": dict.InsufficientArkhoins})
			return
		}
		if errors.Is(err, repository.ErrComponentAlreadyOwned) {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": dict.ComponentAlreadyOwned})
			return
		}
		if errors.Is(err, repository.ErrComponentOwnerCannotBuy) {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": dict.ComponentOwnerCannotBuy})
			return
		}

		log.Print(err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": dict.InternalServerError})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": dict.ComponentBought})
}

func SellComponentHandler(ctx *gin.Context) {
	dict := translations.GetTranslation(ctx)

	issuerID := ctx.Keys["auth_user"].(*dto.UserProfile).ID
	componentID := ctx.Keys["component_lookup"].(*dto.ComponentInfo).ID

	if err := service.Component.Sell(issuerID, componentID); err != nil {
		if errors.Is(err, repository.ErrComponentNotOwned) {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": dict.ComponentNotOwned})
			return
		}
		if errors.Is(err, repository.ErrComponentOwnerCannotSell) {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": dict.ComponentOwnerCannotSell})
			return
		}

		log.Print(err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": dict.InternalServerError})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": dict.ComponentSold})
}

func PublishComponentHandler(ctx *gin.Context) {
	dict := translations.GetTranslation(ctx)

	component := ctx.Keys["component_lookup"].(*dto.ComponentInfo)

	if err := service.Component.Publish(component.ID); err != nil {
		if errors.Is(err, repository.ErrComponentAlreadyPublic) {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": dict.ComponentAlreadyPublic})
			return
		}

		log.Print(err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": dict.InternalServerError})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": dict.ComponentPublished})
}

func UnpublishComponentHandler(ctx *gin.Context) {
	dict := translations.GetTranslation(ctx)

	component := ctx.Keys["component_lookup"].(*dto.ComponentInfo)

	if err := service.Component.Unpublish(component.ID); err != nil {
		if errors.Is(err, repository.ErrComponentNotPublic) {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": dict.ComponentNotPublic})
			return
		}

		log.Print(err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": dict.InternalServerError})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": dict.ComponentUnpublished})
}

func DeleteComponentHandler(ctx *gin.Context) {
	dict := translations.GetTranslation(ctx)

	component := ctx.Keys["component_lookup"].(*dto.ComponentInfo)

	if err := service.Component.SafeDelete(component.ID); err != nil {
		if errors.Is(err, repository.ErrComponentAlreadyTrashed) {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": dict.ComponentAlreadyTrashed})
			return
		}

		log.Print(err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": dict.InternalServerError})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": dict.ComponentTrashed})
}

func DeleteComponenttForceHandler(ctx *gin.Context) {
	dict := translations.GetTranslation(ctx)

	component := ctx.Keys["component_lookup"].(*dto.ComponentInfo)

	if err := service.Component.UnsafeDelete(component.ID); err != nil {
		if errors.Is(err, repository.ErrComponentNotTrashed) {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": dict.ComponentNotTrashed})
			return
		}

		log.Print(err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": dict.InternalServerError})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": dict.ComponentDeleted})
}

func RestoreComponentHandler(ctx *gin.Context) {
	dict := translations.GetTranslation(ctx)

	component := ctx.Keys["component_lookup"].(*dto.ComponentInfo)

	if err := service.Component.Restore(component.ID); err != nil {
		if errors.Is(err, repository.ErrComponentNotTrashed) {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": dict.ComponentNotTrashed})
			return
		}

		log.Print(err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": dict.InternalServerError})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": dict.ComponentRestored})
}
