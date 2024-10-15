package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/swibly/swibly-api/config"
	"github.com/swibly/swibly-api/internal/model/dto"
	"github.com/swibly/swibly-api/translations"
)

func UserPrivacy(requiredShow dto.UserShow) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		dict := translations.GetTranslation(ctx)

		user := ctx.Keys["user_lookup"].(*dto.UserProfile)
		issuer := ctx.Keys["auth_user"].(*dto.UserProfile)

		if user.Username != issuer.Username {
			isAllowed := true

			if requiredShow.Profile == true && !issuer.HasPermissions(config.Permissions.ManageUser) {
				isAllowed = false
			}

			if requiredShow.Image == true && !issuer.HasPermissions(config.Permissions.ManageUser) {
				isAllowed = false
			}

			if requiredShow.Comments == true && !issuer.HasPermissions(config.Permissions.ManageUser) {
				isAllowed = false
			}

			if requiredShow.Favorites == true && !issuer.HasPermissions(config.Permissions.ManageProjects) {
				isAllowed = false
			}

			if requiredShow.Projects == true && !issuer.HasPermissions(config.Permissions.ManageProjects) {
				isAllowed = false
			}

			if requiredShow.Components == true && !issuer.HasPermissions(config.Permissions.ManageProjects) {
				isAllowed = false
			}

			if requiredShow.Followers == true && !issuer.HasPermissions(config.Permissions.ManageUser) {
				isAllowed = false
			}

			if requiredShow.Following == true && !issuer.HasPermissions(config.Permissions.ManageUser) {
				isAllowed = false
			}

			if requiredShow.Inventory == true && !issuer.HasPermissions(config.Permissions.ManageUser) {
				isAllowed = false
			}

			if requiredShow.Formations == true && !issuer.HasPermissions(config.Permissions.ManageUser) {
				isAllowed = false
			}

			if !isAllowed {
				ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": dict.ProjectMissingPermissions})
				return
			}
		}

		ctx.Next()
	}
}
