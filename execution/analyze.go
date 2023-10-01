package execution

import (
	"database/sql"
	"fmt"
	"regexp"
	"strings"
)

func AnalyzeQuery(db *sql.DB, query string) error {
	rows, err := db.Query("EXPLAIN ANALYZE " + query)
	if err != nil {
		fmt.Printf("Failed to analyze query %s: %s\n", query, err)
		return err
	}
	defer rows.Close()

	seqScanCount := 0
	var tables []string
	var line string

	seqScanRegex := regexp.MustCompile(`Seq Scan on (\w+)`) // Regular expression to match "Seq Scan on [table_name]"
	sortRegex := regexp.MustCompile(`Sort[^\n]+`)
	bitmapRegex := regexp.MustCompile(`Bitmap Heap Scan on (\w+)`)
	nestedLoopRegex := regexp.MustCompile(`Nested Loop`)
	hashJoinRegex := regexp.MustCompile(`Hash Join`)
	mergeJoinRegex := regexp.MustCompile(`Merge Join`)

	for rows.Next() {
		if err := rows.Scan(&line); err != nil {
			fmt.Printf("Failed to scan line: %s\n", err)
			return err
		}

		matches := seqScanRegex.FindStringSubmatch(line)
		if len(matches) > 1 {
			seqScanCount++
			tables = append(tables, matches[1]) // The second item in matches will be the table name
		}
		// Check for different operations and provide suggestions
		if seqScanRegex.MatchString(line) {
			fmt.Println("Suggestion: Consider using an Index Scan instead of Seq Scan.")
		}
		if sortRegex.MatchString(line) {
			fmt.Println("Suggestion: An index could be used to avoid sorts.")
		}
		if bitmapRegex.MatchString(line) {
			fmt.Println("Suggestion: See if the query could be rewritten to use an Index Scan.")
		}
		if nestedLoopRegex.MatchString(line) {
			fmt.Println("Suggestion: Analyze whether the tables being joined have appropriate indexes.")
		}
		if hashJoinRegex.MatchString(line) {
			fmt.Println("Hash Join: Generally good for medium-to-large data sets.")
		}
		if mergeJoinRegex.MatchString(line) {
			fmt.Println("Suggestion: Check if an existing index can be used to speed up a Merge Join.")
		}
	}

	if seqScanCount > 0 {
		fmt.Printf("Found %d Seq Scan operations on tables: %s\n", seqScanCount, strings.Join(tables, ", "))
	}

	if err := rows.Err(); err != nil {
		fmt.Printf("Row iteration error: %s\n", err)
		return err
	}

	return nil
}
