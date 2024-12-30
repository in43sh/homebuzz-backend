package database

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

var Db *sql.DB

func ConnectDatabase() {
	host := os.Getenv("DATABASE_HOST")
	port := os.Getenv("DATABASE_PORT")
	user := os.Getenv("DATABASE_USER")
	dbname := os.Getenv("DATABASE_NAME")
	pass := os.Getenv("DATABASE_PASSWORD")

	if host == "" || port == "" || user == "" || dbname == "" || pass == "" {
		fmt.Println("Missing one or more database connection details in environment variables")
		panic("Please check your Render environment settings for missing variables")
	}

	psqlSetup := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=require",
		host, port, user, dbname, pass)

	db, errSql := sql.Open("postgres", psqlSetup)
	if errSql != nil {
		fmt.Println("Error while connecting to the database:", errSql)
		panic(errSql)
	}

	err := db.Ping()
	if err != nil {
		fmt.Println("Error pinging database:", err)
		panic(err)
	} else {
		Db = db
		fmt.Println("Successfully connected to the database!")
	}
}
