package main

import (
	"db-experiments/config"
	"db-experiments/database"
	"fmt"
	"sync"
)

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
		go database.RunQuery(db, query, &wg)
	}
	wg.Wait()

	fmt.Println("Successfully connected to and initialized the Postgres database.")
}
