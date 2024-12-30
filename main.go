package main

import (
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/in43sh/homebuzz-backend/database"
	productRoutes "github.com/in43sh/homebuzz-backend/routes/product" // Correct package import
	userRoutes "github.com/in43sh/homebuzz-backend/routes/user"       // Correct package import
)

func main() {
	route := gin.Default()
	route.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	database.ConnectDatabase()

	route.GET("/ping", func(context *gin.Context) {
		context.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	// Use the AddUser function from the `user` package
	route.GET("/users", userRoutes.GetUsers)
	route.POST("/users", userRoutes.AddUser)
	route.GET("/products", productRoutes.GetProducts)

	err := route.Run(":8080")
	if err != nil {
		panic(err)
	}
}
