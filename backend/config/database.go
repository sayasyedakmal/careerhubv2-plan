package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/microsoft/go-mssqldb"
)

var DB *sql.DB

func InitDB() {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	// Use key-value format instead of URL format — handles named instances
	// (e.g. localhost\careerhubv2) and passwords with special characters correctly.
	dsn := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%s;database=%s",
		host, user, password, port, dbName)

	db, err := sql.Open("sqlserver", dsn)
	if err != nil {
		log.Fatalf("Failed to open database connection: %v", err)
	}

	if err = db.Ping(); err != nil {
		log.Fatalf("Failed to reach SQL Server — check DB_HOST, DB_PORT, DB_USER, DB_PASSWORD, DB_NAME in .env: %v", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)

	log.Println("Connected to SQL Server successfully")
	DB = db
}
