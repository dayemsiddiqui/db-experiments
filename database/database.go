package database

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"os/exec"
	"sync"
	"time"
)

type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

func ConnectDatabase(cfg *DBConfig, retries int) (*sql.DB, error) {
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName)

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

func LoadDatabaseDump(cfg *DBConfig, filepath string) error {
	cmd := exec.Command(
		"bash",
		"-c",
		fmt.Sprintf("PGPASSWORD=%s psql -h %s -p %s -U %s -d %s < %s",
			cfg.Password, cfg.Host, cfg.Port, cfg.User, cfg.DBName, filepath),
	)

	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("Failed to load database dump: %s", err)
		log.Printf("Command output: %s", string(out))
		return err
	}
	log.Printf("Successfully loaded database dump.")
	return nil
}

func RunQuery(db *sql.DB, query string, wg *sync.WaitGroup) {
	defer wg.Done()

	rows, err := db.Query(query)
	if err != nil {
		fmt.Printf("Failed to execute query %s: %s\n", query, err)
		return
	}
	defer rows.Close()

	var count int
	for rows.Next() {
		err := rows.Scan(&count)
		if err != nil {
			fmt.Printf("Failed to read result for query %s: %s\n", query, err)
			return
		}
	}

	fmt.Printf("Query: %s, Result: %d\n", query, count)
}
