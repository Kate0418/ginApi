package main

import (
	"github.com/gin-gonic/gin"
	"back/routes"
)

func main() {
	r := gin.Default()
    routes.V1(r)
    r.Run(":8080")
}
