package main

import (
	"fmt"
	"github.com/alfredoptarigan/go-jwt/controllers"
	"github.com/alfredoptarigan/go-jwt/initializers"
	"github.com/alfredoptarigan/go-jwt/middlewares"
	"github.com/gin-gonic/gin"
)

func init() {
	fmt.Println("Init")
	initializers.LoadEnvVariables()
	initializers.ConnectDatabase()
	initializers.SyncDatabase()
}

func main() {
	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello World",
		})
	})

	r.POST("/sign-up", controllers.Signup)
	r.POST("/sign-in", controllers.Login)
	r.GET("/validate", middlewares.AuthProtection, controllers.Validate)

	r.Run()
}
