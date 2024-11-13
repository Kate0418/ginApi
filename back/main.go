package main

import (
	"github.com/gin-gonic/gin"
	"back/routes"
)

func main() {
	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
        c.JSON(200, gin.H{
            "message": "Hello Gin!",
        })
    })
    routes.V1(r)
    r.Run(":8080")
}
