package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/in43sh/homebuzz-backend/database"
	_ "github.com/in43sh/homebuzz-backend/docs"
	productRoutes "github.com/in43sh/homebuzz-backend/routes/product"
	userRoutes "github.com/in43sh/homebuzz-backend/routes/user"
	swaggerFiles "github.com/swaggo/files" // swagger embed files
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Your API
// @version 1.0
// @description This is a sample server for managing authentication.
// @host localhost:8080
// @BasePath /

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

	// Swagger endpoint
	route.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// User routes
	route.POST("/register", userRoutes.Register)
	route.POST("/login", userRoutes.Login)
	route.GET("/users", userRoutes.GetUsers)
	route.GET("/users/:id", userRoutes.GetUser)
	route.DELETE("/users/:id", userRoutes.DeleteUser)

	// Product routes
	route.POST("/products", productRoutes.AddProduct)
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
