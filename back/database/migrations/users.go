package migrations

import (
	"back/api/models"
	"back/database"
)

func Users() {
	database.Gorm().AutoMigrate(&models.User{})
}
