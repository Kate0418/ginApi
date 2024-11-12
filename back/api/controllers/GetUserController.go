package controllers

import (
    "github.com/gin-gonic/gin"
    "back/database"
    "back/api/models"
    "net/http"
)

func GetUserController(c *gin.Context) {
	db := database.gorm()
}
