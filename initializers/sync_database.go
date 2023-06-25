package initializers

import "github.com/alfredoptarigan/go-jwt/models"

func SyncDatabase() {
	DB.AutoMigrate(&models.User{})
}
