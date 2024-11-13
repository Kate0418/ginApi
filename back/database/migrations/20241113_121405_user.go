package migrations

import (
    "back/api/models"
    "back/database"
)

func CreateUser() {
    db := database.Gorm()
    db.AutoMigrate(&models.User{})
}
