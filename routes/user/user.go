package routes

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/in43sh/homebuzz-backend/database"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       int64  `bun:",pk,autoincrement"`
	Username string `bun:"username,unique,notnull" json:"username" binding:"required" example:"johndoe"`
	Password string `bun:"password,notnull" json:"password" binding:"required" example:"password123"`
}

type SuccessResponse struct {
	Message string `json:"message" example:"User successfully created"`
}

type ErrorResponse struct {
	Error string `json:"error" example:"Invalid input"`
}

var jwtKey = []byte("your_secret_key")

type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// @Summary Register a new user
// @Description Register a new user by providing username and password
// @Tags Auth
// @Accept  json
// @Produce  json
// @Param user body User true "User credentials"
// @Success 200 {object} SuccessResponse "User successfully created"
// @Failure 400 {object} ErrorResponse "Invalid input"
// @Failure 409 {object} ErrorResponse "User already exists"
// @Failure 500 {object} ErrorResponse "Failed to create user"
// @Router /register [post]
func Register(ctx *gin.Context) {
	var user User

	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid input"})
		return
	}

	existingUser := new(User)
	err := database.BunDB.NewSelect().
		Model(existingUser).
		Where("username = ?", user.Username).
		Scan(context.Background())
	if err == nil {
		ctx.AbortWithStatusJSON(http.StatusConflict, ErrorResponse{Error: "User already exists"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to encrypt password"})
		return
	}
	user.Password = string(hashedPassword)

	_, err = database.BunDB.NewInsert().Model(&user).Exec(context.Background())
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to create user"})
		return
	}

	ctx.JSON(http.StatusOK, SuccessResponse{Message: "User successfully created"})
}

// @Summary Login user
// @Description Authenticate a user and return a JWT token.
// @Tags Auth
// @Accept  json
// @Produce  json
// @Param credentials body User true "User credentials"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /login [post]
func Login(ctx *gin.Context) {
	var credentials User

	if err := ctx.ShouldBindJSON(&credentials); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	storedUser := new(User)
	err := database.BunDB.NewSelect().
		Model(storedUser).
		Where("username = ?", credentials.Username).
		Scan(context.Background())
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(storedUser.Password), []byte(credentials.Password))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		Username: credentials.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message":  "Login successful",
		"username": credentials.Username,
		"token":    tokenString,
	})
}

// @Summary Get all users
// @Description Retrieve a list of all users in the system
// @Tags Users
// @Accept  json
// @Produce  json
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /users [get]
func GetUsers(ctx *gin.Context) {
	var users []User

	err := database.BunDB.NewSelect().
		Model(&users).
		Scan(context.Background())
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Couldn't fetch users"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"users": users})
}

// @Summary Get a specific user by ID
// @Description Retrieve a single user by their ID
// @Tags Users
// @Accept  json
// @Produce  json
// @Param id path int64 true "User ID"
// @Success 200 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /users/{id} [get]
func GetUser(ctx *gin.Context) {
	id := ctx.Param("id")

	user := new(User)
	err := database.BunDB.NewSelect().
		Model(user).
		Where("id = ?", id).
		Scan(context.Background())
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"user": user})
}

// @Summary Delete a user by ID
// @Description Delete a user from the system by their ID
// @Tags Users
// @Accept  json
// @Produce  json
// @Param id path int64 true "User ID"
// @Success 200 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /users/{id} [delete]
func DeleteUser(ctx *gin.Context) {
	id := ctx.Param("id")

	result, err := database.BunDB.NewDelete().
		Model((*User)(nil)).
		Where("id = ?", id).
		Exec(context.Background())
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "User successfully deleted"})
}
