package reporting

import (
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"log"
)

type Report struct {
	PgStatStatements    []map[string]interface{} `json:"pg_stat_statements"`
	PgStatioUserIndexes []map[string]interface{} `json:"pg_statio_user_indexes"`
}

func GenerateReport(db *sql.DB) {
	report := Report{}

	// Execute the first query
	rows, err := db.Query("SELECT * from pg_stat_statements;")
	if err != nil {
		log.Fatal(err)
	}
	report.PgStatStatements = parseRows(rows)
	rows.Close()

	// Execute the second query
	rows, err = db.Query("SELECT * FROM pg_statio_user_indexes;")
	if err != nil {
		log.Fatal(err)
	}
	report.PgStatioUserIndexes = parseRows(rows)
	rows.Close()

	// Serialize the report to JSON
	reportJSON, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	// Write the JSON to a file
	err = ioutil.WriteFile("reports.json", reportJSON, 0644)
	if err != nil {
		log.Fatal(err)
	}
}

func parseRows(rows *sql.Rows) []map[string]interface{} {
	columns, err := rows.Columns()
	if err != nil {
		log.Fatal(err)
	}

	var result []map[string]interface{}
	for rows.Next() {
		values := make([]interface{}, len(columns))
		pointers := make([]interface{}, len(columns))
		for i := range columns {
			pointers[i] = &values[i]
		}

		if err := rows.Scan(pointers...); err != nil {
			log.Fatal(err)
		}

		row := make(map[string]interface{})
		for i, colName := range columns {
			val := values[i]
			b, ok := val.([]byte)
			if ok {
				row[colName] = string(b)
			} else {
				row[colName] = val
			}
		}

		result = append(result, row)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	return result
}
