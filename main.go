package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/in43sh/homebuzz-backend/database"
	productRoutes "github.com/in43sh/homebuzz-backend/routes/product"
	userRoutes "github.com/in43sh/homebuzz-backend/routes/user"
)

func main() {
	route := gin.Default()

	allowOrigins := []string{"http://localhost:3000"}
	if os.Getenv("GIN_MODE") == "release" {
		allowOrigins = []string{"https://homebuzz-backend.onrender.com", "https://homebuzz.netlify.app"}
	}
	fmt.Printf("allowOrigins: %s\n", allowOrigins)

	route.Use(cors.New(cors.Config{
		AllowOrigins:     allowOrigins,
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

	// User routes
	route.POST("/register", userRoutes.Register)
	route.POST("/login", userRoutes.Login)
	route.GET("/users", userRoutes.GetUsers)
	route.GET("/users/:username", userRoutes.GetUserByUsername)
	route.DELETE("/users/:username", userRoutes.DeleteUserByUsername)

	// Product routes
	route.GET("/products", productRoutes.GetProducts)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("Server is running on port %s\n", port)

	err := route.Run(":" + port)
	if err != nil {
		panic(err)
	}
}
