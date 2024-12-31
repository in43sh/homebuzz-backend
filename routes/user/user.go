package routes

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/in43sh/homebuzz-backend/database"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
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

	var existingUser string
	err := database.Db.QueryRow("SELECT username FROM users WHERE username = $1", user.Username).Scan(&existingUser)
	if err == nil {
		ctx.AbortWithStatusJSON(http.StatusConflict, gin.H{"error": "User already exists"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to encrypt password"})
		return
	}

	_, err = database.Db.Exec("INSERT INTO users(username, password) VALUES ($1, $2)", user.Username, string(hashedPassword))
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

	var storedPassword string

	err := database.Db.QueryRow("SELECT password FROM users WHERE username = $1", credentials.Username).Scan(&storedPassword)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(credentials.Password))
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
	rows, err := database.Db.Query("SELECT username, password FROM users")
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Couldn't fetch users"})
		return
	}
	defer rows.Close()

	users := []User{}

	for rows.Next() {
		var user User
		if err := rows.Scan(&user.Username, &user.Password); err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Error reading user data"})
			return
		}
		users = append(users, user)
	}

	ctx.JSON(http.StatusOK, gin.H{"users": users})
}

func GetUserByUsername(ctx *gin.Context) {
	username := ctx.Param("username")

	var user User
	err := database.Db.QueryRow("SELECT username, password FROM users WHERE username = $1", username).Scan(&user.Username, &user.Password)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"user": user})
}

func DeleteUserByUsername(ctx *gin.Context) {
	username := ctx.Param("username")

	result, err := database.Db.Exec("DELETE FROM users WHERE username = $1", username)
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
