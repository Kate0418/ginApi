package routes

import (
    "github.com/gin-gonic/gin"
    "back/api/controllers"
)

func V1(r *gin.Engine) {
	v1 := r.Group("/v1")
	{
		v1.GET("/users", controllers.GetUser)
	}
}
