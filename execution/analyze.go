package execution

import (
	"database/sql"
	"fmt"
)

// AnalyzeQuery runs EXPLAIN ANALYZE and applies analysis rules
func AnalyzeQuery(db *sql.DB, query string) (string, error) {
	rows, err := db.Query("EXPLAIN ANALYZE " + query)
	if err != nil {
		fmt.Printf("Failed to analyze query %s: %s\n", query, err)
		return "", err
	}
	defer rows.Close()

	var line string
	var tables []string

	var output string

	for rows.Next() {
		if err := rows.Scan(&line); err != nil {
			fmt.Printf("Failed to scan line: %s\n", err)
			return "", err
		}

		output += line + "\n"
	}

	if len(tables) > 0 {
		//fmt.Printf("Tables involved: %s\n", tables)
	}

	if err := rows.Err(); err != nil {
		fmt.Printf("Row iteration error: %s\n", err)
		return "", err
	}

	return "", nil
}

func extractTable(matches []string) []string {
	if len(matches) > 1 {
		return matches[1:]
	}
	return nil
}
