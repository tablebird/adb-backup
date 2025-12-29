package admin

import (
	"adb-backup/internal/config"
	"adb-backup/internal/database"
)

func InitAdmin() {

	count, err := database.CountLocalUsers()
	if err != nil {
		// Handle error
		return
	}

	conf := config.Web
	if count == 0 {
		adminUser := &database.User{
			Name:     conf.AdminName,
			Password: conf.AdminPassword,
		}
		database.CreateUser(adminUser)
	}

}
