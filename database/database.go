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
		fmt.Println("Error loading .env file")
	}

	host := os.Getenv("DATABASE_HOST")
	port := os.Getenv("DATABASE_PORT")
	if port == "" {
		port = "5432"
	}
	user := os.Getenv("DATABASE_USER")
	dbname := os.Getenv("DATABASE_NAME")
	pass := os.Getenv("DATABASE_PASSWORD")

	ginMode := os.Getenv("GIN_MODE")
	if ginMode == "" {
		ginMode = "debug"
	}

	sslMode := "disable"
	if ginMode == "release" {
		sslMode = "require"
	}

	fmt.Printf("DATABASE_HOST: %s\n", host)
	fmt.Printf("DATABASE_PORT: %s\n", port)
	fmt.Printf("DATABASE_USER: %s\n", user)
	fmt.Printf("DATABASE_NAME: %s\n", dbname)
	fmt.Printf("DATABASE_PASSWORD: %s\n", pass)
	fmt.Printf("SSLMODE: %s\n", sslMode)
	fmt.Printf("GYN_MODE: %s\n", ginMode)

	psqlSetup := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		host, port, user, dbname, pass, sslMode)

	var errSql error
	Db, errSql = sql.Open("postgres", psqlSetup)
	if errSql != nil {
		fmt.Println("Error while connecting to the database:", errSql)
		panic(errSql)
	}

	err = Db.Ping()
	if err != nil {
		fmt.Println("Error pinging database:", err)
		panic(err)
	}

	Db.SetMaxOpenConns(10)
	Db.SetMaxIdleConns(5)
	Db.SetConnMaxLifetime(0)

	fmt.Println("Successfully connected to the database!")
}
