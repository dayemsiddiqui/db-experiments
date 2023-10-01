package execution

import (
	"database/sql"
	"fmt"
	"regexp"
)

// Rule represents a pattern and corresponding suggestion
type Rule struct {
	pattern *regexp.Regexp
	suggest string
	extract func(matches []string) []string
}

// AnalyzeQuery runs EXPLAIN ANALYZE and applies analysis rules
func AnalyzeQuery(db *sql.DB, query string) error {
	rows, err := db.Query("EXPLAIN ANALYZE " + query)
	if err != nil {
		fmt.Printf("Failed to analyze query %s: %s\n", query, err)
		return err
	}
	defer rows.Close()

	rules := []Rule{
		{regexp.MustCompile(`Seq Scan on (\w+)`), "Suggestion: Consider using an Index Scan instead of Seq Scan.", extractTable},
		{regexp.MustCompile(`Sort[^\n]+`), "Suggestion: An index could be used to avoid sorts.", nil},
		{regexp.MustCompile(`Bitmap Heap Scan on (\w+)`), "Suggestion: See if the query could be rewritten to use an Index Scan.", extractTable},
		{regexp.MustCompile(`Nested Loop`), "Suggestion: Analyze whether the tables being joined have appropriate indexes.", nil},
		{regexp.MustCompile(`Hash Join`), "Hash Join: Generally good for medium-to-large data sets.", nil},
		{regexp.MustCompile(`Merge Join`), "Suggestion: Check if an existing index can be used to speed up a Merge Join.", nil},
	}

	var line string
	var tables []string

	for rows.Next() {
		if err := rows.Scan(&line); err != nil {
			fmt.Printf("Failed to scan line: %s\n", err)
			return err
		}

		for _, rule := range rules {
			matches := rule.pattern.FindStringSubmatch(line)
			if len(matches) > 0 {
				//fmt.Println(rule.suggest)
				if rule.extract != nil {
					tables = append(tables, rule.extract(matches)...)
				}
			}
		}
	}

	if len(tables) > 0 {
		//fmt.Printf("Tables involved: %s\n", tables)
	}

	if err := rows.Err(); err != nil {
		fmt.Printf("Row iteration error: %s\n", err)
		return err
	}

	return nil
}

func extractTable(matches []string) []string {
	if len(matches) > 1 {
		return matches[1:]
	}
	return nil
}
