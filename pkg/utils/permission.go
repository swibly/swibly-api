package utils

import (
	"slices"

	"github.com/swibly/swibly-api/config"
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
