package middleware

import (
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/swibly/swibly-api/config"
	"github.com/swibly/swibly-api/internal/model/dto"
	"github.com/swibly/swibly-api/internal/service"
	"github.com/swibly/swibly-api/translations"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func ComponentLookup(ctx *gin.Context) {
	dict := translations.GetTranslation(ctx)

	issuer := ctx.Keys["auth_user"].(*dto.UserProfile)

	componentID, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		log.Print(err)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": dict.ComponentInvalid})
		return
	}

	component, err := service.Component.GetByID(issuer.ID, uint(componentID))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": dict.ComponentNotFound})
			return
		}

		log.Print(err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": dict.InternalServerError})
		return
	}

	if !component.IsPublic && issuer.ID != component.OwnerID {
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": dict.ComponentNotFound})
		return
	}

	ctx.Set("component_lookup", component)
	ctx.Next()
}

// middleware.ComponentLookup must be called before this
func ComponentOwnership(ctx *gin.Context) {
	dict := translations.GetTranslation(ctx)

	issuer := ctx.Keys["auth_user"].(*dto.UserProfile)

	component := ctx.Keys["component_lookup"].(*dto.ComponentInfo)
	if issuer.HasPermissions(config.Permissions.ManageStore) || component.OwnerUsername != issuer.Username {
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": dict.ComponentNotFound})
		return
	}

	ctx.Next()
}
