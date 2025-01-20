package main

import (
	"API/database"
	"API/server"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	// Connect to the database
	err := database.InitDB()
	if err != nil {
		panic(err)
	}
	// Start and initialize the server
	server.StartAndInitializeServer()
}
