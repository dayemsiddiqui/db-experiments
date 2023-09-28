package database

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"os/exec"
	"time"
)

func ConnectDatabase(retries int) (*sql.DB, error) {
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

func LoadDatabaseDump(filepath string) error {
	// Replace the following placeholders with actual values
	host := "localhost"
	port := "5432"
	user := "user"
	password := "password"
	dbname := "mydatabase"

	// Set the command to run psql and load the dump
	cmd := exec.Command(
		"bash",
		"-c",
		fmt.Sprintf("PGPASSWORD=%s psql -h %s -p %s -U %s -d %s < %s", password, host, port, user, dbname, filepath),
	)

	// Execute the command
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("Failed to load database dump: %s", err)
		log.Printf("Command output: %s", string(out))
		return err
	}
	log.Printf("Successfully loaded database dump.")
	return nil
}
