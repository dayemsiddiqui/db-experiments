package main

import (
	"db-experiments/config"
	"db-experiments/database"
	"db-experiments/docker"
	"db-experiments/execution"
	"db-experiments/reporting"
	"fmt"
	"sync"
)

func main() {
	docker.RunContainer()
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
	db.SetMaxOpenConns(20)
	defer db.Close()
	// Load dump file from config
	err = database.LoadDatabaseDump(dbCfg, cfg.DumpPath)
	if err != nil {
		fmt.Printf("Failed to load the database dump: %s\n", err)
		return
	}

	// Run queries from config in separate goroutines based on their traffic percent
	var wg sync.WaitGroup
	for _, queryConfig := range cfg.Queries {
		wg.Add(1)
		runQueryConfig := execution.RunQueryConfig{
			QueryConfig: queryConfig,
			Config:      cfg,
		}
		go execution.RunQuery(db, runQueryConfig, &wg)
	}
	wg.Wait()

	fmt.Println("Successfully connected to and initialized the Postgres database.")

	fmt.Println("Generating reports...")
	// Generate reports
	reporting.GenerateReport(db)
}
