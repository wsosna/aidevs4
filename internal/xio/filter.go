package xio

import (
	"encoding/csv"
	"strconv"
	"strings"
)

// FilterCsvFromString parses CSV content, keeps rows where birthDate year is between from and to (inclusive),
// and optionally filters by other columns positionally (empty string = skip that column).
func FilterCsvFromString(content string, from, to int, columnFilters ...string) ([][]string, error) {
	records, err := csv.NewReader(strings.NewReader(content)).ReadAll()
	if err != nil {
		return nil, err
	}
	if len(records) == 0 {
		return nil, nil
	}

	header := records[0]
	birthDateIdx := -1
	for i, h := range header {
		if h == "birthDate" {
			birthDateIdx = i
			break
		}
	}

	result := [][]string{header}
	for _, row := range records[1:] {
		if birthDateIdx >= 0 && birthDateIdx < len(row) {
			year, err := strconv.Atoi(strings.SplitN(row[birthDateIdx], "-", 2)[0])
			if err != nil || year < from || year > to {
				continue
			}
		}

		match := true
		for i, filter := range columnFilters {
			if filter == "" || i >= len(row) {
				continue
			}
			if row[i] != filter {
				match = false
				break
			}
		}
		if match {
			result = append(result, row)
		}
	}

	return result, nil
}
