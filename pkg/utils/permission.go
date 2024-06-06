package utils

import (
	"slices"

	"github.com/devkcud/arkhon-foundation/arkhon-api/config"
)

func HasPermissions(list []string, lookup ...string) bool {
	if slices.Contains(list, config.Permissions.Admin) {
		return true
	}

	return true
}
