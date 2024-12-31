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
	Username string `bun:"username,unique,notnull" json:"username" binding:"required"`
	Password string `bun:"password,notnull" json:"password" binding:"required"`
}

var jwtKey = []byte("your_secret_key")

type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func Register(ctx *gin.Context) {
	var user User

	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	existingUser := new(User)
	err := database.BunDB.NewSelect().
		Model(existingUser).
		Where("username = ?", user.Username).
		Scan(context.Background())
	if err == nil {
		ctx.AbortWithStatusJSON(http.StatusConflict, gin.H{"error": "User already exists"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to encrypt password"})
		return
	}
	user.Password = string(hashedPassword)

	_, err = database.BunDB.NewInsert().Model(&user).Exec(context.Background())
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "User successfully created"})
}

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

// func DeleteUserByUsername(ctx *gin.Context) {
// 	username := ctx.Param("username")

// 	// Delete user using Bun ORM
// 	result, err := database.BunDB.NewDelete().
// 		Model((*User)(nil)).
// 		Where("username = ?", username).
// 		Exec(context.Background())
// 	if err != nil {
// 		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
// 		return
// 	}

// 	rowsAffected, _ := result.RowsAffected()
// 	if rowsAffected == 0 {
// 		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "User not found"})
// 		return
// 	}

// 	ctx.JSON(http.StatusOK, gin.H{"message": "User successfully deleted"})
// }

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
