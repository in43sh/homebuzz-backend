package routes

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/in43sh/homebuzz-backend/database"
)

// User represents the structure of a user
type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// AddUser adds a new user to the database
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
