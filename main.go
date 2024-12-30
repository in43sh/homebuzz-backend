package main

import (
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
	// Create a new Gin router
	route := gin.Default()

	// Configure CORS
	route.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Initialize the database connection
	database.ConnectDatabase()

	// Define routes
	route.GET("/ping", func(context *gin.Context) {
		context.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	// User routes
	route.GET("/users", userRoutes.GetUsers)
	route.POST("/users", userRoutes.AddUser)

	// Product routes
	route.GET("/products", productRoutes.GetProducts)

	// Get the port from the environment variable
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default to port 8080 if PORT is not set
	}

	// Start the server
	err := route.Run(":" + port)
	if err != nil {
		panic(err)
	}
}
