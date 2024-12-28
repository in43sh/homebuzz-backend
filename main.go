package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/in43sh/homebuzz-backend/database"
	userRoutes "github.com/in43sh/homebuzz-backend/routes/user" // Correct package import
)

func main() {
	route := gin.Default()
	database.ConnectDatabase()

	route.GET("/ping", func(context *gin.Context) {
		context.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	// Use the AddUser function from the `user` package
	route.POST("/add", userRoutes.AddUser)

	err := route.Run(":8080")
	if err != nil {
		panic(err)
	}
}
