package controllers

import (
    "github.com/gin-gonic/gin"
    "back/database"
    "back/api/models"
    "net/http"
)

func GetUser(c *gin.Context) {
	var users []models.User
    db := database.Gorm()

    result := db.Find(&users)
    if result.Error != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": "Failed to fetch users",
        })
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "data": users,
    })
}
