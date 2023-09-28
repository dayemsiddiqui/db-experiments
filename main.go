package main

import (
	"db-experiments/database"
	"db-experiments/docker"
	"fmt"
	"log"
)

func main() {
	docker.RunContainer()
	db, err := database.ConnectDatabase(5) // Try connecting 5 times
	if err != nil {
		log.Fatalf("Failed to connect to the database: %s", err)
	}
	err = database.LoadDatabaseDump("./dump.sql")
	if err != nil {
		log.Fatalf("Failed to load the database dump: %s", err)
	}
	defer db.Close()

	fmt.Println("Successfully connected to the Postgres database.")
}
