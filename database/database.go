package database

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

var BunDB *bun.DB

func ConnectDatabase() {
	// Load environment variables
	if os.Getenv("GIN_MODE") != "release" {
		err := godotenv.Load()
		if err != nil {
			fmt.Println("Error loading .env file")
		}
	}

	// Environment variables
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

	// Bun database connection
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		user, pass, host, port, dbname, sslMode)

	pgConn := pgdriver.NewConnector(pgdriver.WithDSN(dsn))
	sqlDB := sql.OpenDB(pgConn) // Use sql.OpenDB with pgdriver.Connector

	// Validate connection
	err := sqlDB.Ping()
	if err != nil {
		fmt.Println("Error pinging database:", err)
		panic(err)
	}

	// Configure connection pool
	sqlDB.SetMaxOpenConns(10)
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// Initialize Bun
	BunDB = bun.NewDB(sqlDB, pgdialect.New())

	fmt.Println("Successfully connected to the database with Bun!")
}
