package utils

import (
	"slices"

	"github.com/devkcud/arkhon-foundation/arkhon-api/config"
	"github.com/gin-gonic/gin"
)

// There is no need to pass anything to check for admin roles
func HasPermissions(list []string, lookup ...string) bool {
	if slices.Contains(list, config.Permissions.Admin) {
		return true
	}

	if len(lookup) == 0 {
		return false
	}

	for _, perm := range lookup {
		if !slices.Contains(list, perm) {
			return false
		}
	}

	return true
}

// WARN: must have the GetPermissionsMiddleware in the route before using
func HasPermissionsByContext(ctx *gin.Context, lookup ...string) bool {
	return HasPermissions(ctx.Keys["permissions"].([]string), lookup...)
}
