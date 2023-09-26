package main

import (
	"database/sql"
	"fmt"
	"log"
	"os/exec"
	"strings"
	"time"

	_ "github.com/lib/pq"
)

func isContainerRunning() bool {
	cmd := exec.Command("docker-compose", "ps", "-q", "postgres")
	output, err := cmd.Output()
	if err != nil {
		log.Fatalf("Failed to check if the container is running: %s", err)
	}
	return strings.TrimSpace(string(output)) != ""
}

func connectDatabase(retries int) (*sql.DB, error) {
	connStr := "host=localhost port=5432 user=user password=password dbname=mydatabase sslmode=disable"
	var db *sql.DB
	var err error

	for i := 0; i < retries; i++ {
		db, err = sql.Open("postgres", connStr)
		if err != nil {
			log.Printf("Failed to connect to database. Attempt %d/%d: %s", i+1, retries, err)
			time.Sleep(5 * time.Second)
			continue
		}

		err = db.Ping()
		if err != nil {
			log.Printf("Failed to ping the database. Attempt %d/%d: %s", i+1, retries, err)
			time.Sleep(5 * time.Second)
			continue
		}

		return db, nil
	}

	return nil, fmt.Errorf("could not connect to the database after %d attempts", retries)
}

func main() {
	if !isContainerRunning() {
		cmd := exec.Command("docker-compose", "up", "-d")
		err := cmd.Run()
		if err != nil {
			log.Fatalf("Failed to start the Postgres container: %s", err)
		}
	}

	db, err := connectDatabase(5) // Try connecting 5 times
	if err != nil {
		log.Fatalf("Failed to connect to the database: %s", err)
	}
	defer db.Close()

	fmt.Println("Successfully connected to the Postgres database.")
}
