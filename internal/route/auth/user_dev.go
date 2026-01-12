//go:build dev

package auth

import (
	"adb-backup/internal/config"
)

func appendUser(h map[string]any) {
	h["User"] = config.Web.AdminName
	h["Password"] = config.Web.AdminPassword
}
