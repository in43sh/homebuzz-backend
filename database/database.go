package database

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var Db *sql.DB

func ConnectDatabase() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error occurred loading .env file, please check")
		panic(err)
	}

	host := os.Getenv("HOST")
	port := os.Getenv("PORT")
	user := os.Getenv("USER")
	dbname := os.Getenv("DB_NAME")
	pass := os.Getenv("PASSWORD")

	if host == "" || port == "" || user == "" || dbname == "" || pass == "" {
		fmt.Println("Missing one or more database connection details in .env file")
		panic("Please check your .env file for missing variables")
	}

	psqlSetup := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=require",
		host, port, user, dbname, pass)

	db, errSql := sql.Open("postgres", psqlSetup)
	if errSql != nil {
		fmt.Println("Error while connecting to the database:", errSql)
		panic(errSql)
	}

	err = db.Ping()
	if err != nil {
		fmt.Println("Error pinging database:", err)
		panic(err)
	} else {
		Db = db
		fmt.Println("Successfully connected to the database!")
	}
}
