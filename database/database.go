package database

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
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
