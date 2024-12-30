package routes

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/in43sh/homebuzz-backend/database"
)

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func AddUser(ctx *gin.Context) {
	body := User{}
	data, err := ctx.GetRawData()
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "User is not defined"})
		return
	}
	err = json.Unmarshal(data, &body)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Bad Input"})
		return
	}

	_, err = database.Db.Exec("INSERT INTO users(username, password) VALUES ($1, $2)", body.Username, body.Password)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Couldn't create the new user."})
	} else {
		ctx.JSON(http.StatusOK, gin.H{"message": "User successfully created."})
	}
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
