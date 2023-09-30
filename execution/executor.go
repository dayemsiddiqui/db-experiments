package execution

import (
	"database/sql"
	"db-experiments/query-preparation"
	"fmt"
	"sync"
)

func RunQuery(db *sql.DB, runConfig RunQueryConfig, wg *sync.WaitGroup) {
	defer wg.Done()
	queryConfig := runConfig.QueryConfig
	cfg := runConfig.Config
	count := int(float64(cfg.Traffic) * queryConfig.TrafficPercent)
	fmt.Println("Running query: ", queryConfig.Name, " ", count, " times")
	query := query_preparation.PrepareQuery(cfg, &queryConfig)
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
