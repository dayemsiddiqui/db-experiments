package execution

import (
	"database/sql"
	"db-experiments/config"
	"db-experiments/query-preparation"
	"fmt"
	"sync"
)

func RunQuery(db *sql.DB, runConfig RunQueryConfig, wg *sync.WaitGroup) {
	defer wg.Done()
	queryConfig := runConfig.QueryConfig
	cfg := runConfig.Config
	count := getQueryCount(cfg, queryConfig)
	fmt.Println("Running query: ", queryConfig.Name, " ", count, " times")
	query := query_preparation.PrepareQuery(cfg, &queryConfig)

	if _, analysisErr := AnalyzeQuery(db, query); analysisErr != nil {
		return
	}

	for i := 0; i < count; i++ {
		rows, err := db.Query(query)
		if err != nil {
			fmt.Printf("Failed to execute query %s: %s\n", query, err)
			return
		}
		err = rows.Close()
		if err != nil {
			return
		}
	}
}

func getQueryCount(cfg *config.Config, queryConfig config.QueryConfig) int {
	count := int(float64(cfg.Traffic) * (float64(queryConfig.TrafficPercent) / 100.0))
	return count
}
