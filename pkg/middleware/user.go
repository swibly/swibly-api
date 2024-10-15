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

			if requiredShow.Profile == true && !user.Show.Profile && !issuer.HasPermissions(config.Permissions.ManageUser) {
				isAllowed = false
			}

			if requiredShow.Image == true && !user.Show.Image && !issuer.HasPermissions(config.Permissions.ManageUser) {
				isAllowed = false
			}

			if requiredShow.Comments == true && !user.Show.Comments && !issuer.HasPermissions(config.Permissions.ManageUser) {
				isAllowed = false
			}

			if requiredShow.Favorites == true && !user.Show.Favorites && !issuer.HasPermissions(config.Permissions.ManageProjects) {
				isAllowed = false
			}

			if requiredShow.Projects == true && !user.Show.Favorites && !issuer.HasPermissions(config.Permissions.ManageProjects) {
				isAllowed = false
			}

			if requiredShow.Components == true && !user.Show.Components && !issuer.HasPermissions(config.Permissions.ManageProjects) {
				isAllowed = false
			}

			if requiredShow.Followers == true && !user.Show.Followers && !issuer.HasPermissions(config.Permissions.ManageUser) {
				isAllowed = false
			}

			if requiredShow.Following == true && !user.Show.Following && !issuer.HasPermissions(config.Permissions.ManageUser) {
				isAllowed = false
			}

			if requiredShow.Inventory == true && !user.Show.Inventory && !issuer.HasPermissions(config.Permissions.ManageUser) {
				isAllowed = false
			}

			if requiredShow.Formations == true && !user.Show.Formations && !issuer.HasPermissions(config.Permissions.ManageUser) {
				isAllowed = false
			}

			if !isAllowed {
				ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": dict.UserMissingPermissions})
				return
			}
		}

		ctx.Next()
	}
}
