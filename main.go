package main

import (
	"API/database"
	"fmt"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func main() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Retrieve variables from environment
	dbUser := os.Getenv("DB_USER")
	dbPwd := os.Getenv("DB_PASSWORD")
	dbTCPHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	// Connect to the database
	err = database.InitDB(dbUser, dbPwd, dbTCPHost, dbPort, dbName)
	if err != nil {
		panic(err)
	}
	fmt.Println(database.RegisterUserFromEmail("test@test.com", "test", "123456"))

	fmt.Println(database.GetTripByID(1))
}
