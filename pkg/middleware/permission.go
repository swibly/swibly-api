package middleware

import (
	"github.com/devkcud/arkhon-foundation/arkhon-api/internal/model/dto"
	"github.com/devkcud/arkhon-foundation/arkhon-api/internal/service"
	"github.com/gin-gonic/gin"
)

// GetPermissionsMiddleware must be after OptionalAuthMiddleware or AuthMiddleware
func GetPermissionsMiddleware(ctx *gin.Context) {
	var issuer *dto.ProfileSearch = nil
	p, exists := ctx.Get("auth_user")

	if !exists {
		ctx.Set("permissions", []string{}) // Set to an empty string so it can be queried anyway
		ctx.Next()
		return
	}

	issuer = p.(*dto.ProfileSearch)

	permissions, err := service.Permission.GetPermissions(issuer.ID)
	if err != nil {
		ctx.Set("permissions", []string{})
		ctx.Next()
		return
	}

	var list []string

	for _, permission := range permissions {
		list = append(list, permission.Name)
	}

	ctx.Set("permissions", list)
}
