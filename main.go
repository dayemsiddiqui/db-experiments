package main

import (
	"database/sql"
	"db-experiments/config"
	"db-experiments/database"
	"fmt"
	"sync"
)

func runQuery(db *sql.DB, query string, wg *sync.WaitGroup) {
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

func main() {
	// Read config from experiments.yaml
	cfg, err := config.ReadConfig()
	if err != nil {
		fmt.Printf("Failed to read config: %s\n", err)
		return
	}

	dbCfg := &database.DBConfig{
		Host:     "localhost",
		Port:     "5432",
		User:     "user",
		Password: "password",
		DBName:   "mydatabase",
	}

	db, err := database.ConnectDatabase(dbCfg, 5) // Try connecting 5 times
	if err != nil {
		fmt.Printf("Failed to connect to the database: %s\n", err)
		return
	}
	defer db.Close()
	// Load dump file from config
	err = database.LoadDatabaseDump(dbCfg, cfg.DumpPath)
	if err != nil {
		fmt.Printf("Failed to load the database dump: %s\n", err)
		return
	}

	// Run queries from config in separate goroutines
	var wg sync.WaitGroup
	for _, query := range cfg.Queries {
		wg.Add(1)
		go runQuery(db, query, &wg)
	}
	wg.Wait()

	fmt.Println("Successfully connected to and initialized the Postgres database.")
}
