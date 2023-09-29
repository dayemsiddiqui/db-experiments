package database

import (
	"database/sql"
	"db-experiments/config"
	_ "db-experiments/config"
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
	// Existing pg_dump code
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

	// Database connection string
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName)

	// Connect to the database
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Printf("Failed to connect to the database: %s", err)
		return err
	}
	defer db.Close()

	// Enable pg_stat_statements if not already enabled
	_, err = db.Exec("CREATE EXTENSION IF NOT EXISTS pg_stat_statements;")
	if err != nil {
		log.Printf("Failed to enable pg_stat_statements: %s", err)
		return err
	}

	// Reset any pre-existing stats
	_, err = db.Exec("SELECT pg_stat_statements_reset();")
	if err != nil {
		log.Printf("Failed to reset pg_stat_statements stats: %s", err)
		return err
	}

	return nil
}

func RunQuery(db *sql.DB, queryConfig config.QueryConfig, count int, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Println("Running query: ", queryConfig.Name, " ", count, " times")
	for i := 0; i < count; i++ {

		rows, err := db.Query(queryConfig.Query)
		if err != nil {
			fmt.Printf("Failed to execute query %s: %s\n", queryConfig.Query, err)
			return
		}
		err = rows.Close()
		if err != nil {
			return
		}
	}
}
