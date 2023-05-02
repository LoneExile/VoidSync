package utils

import (
	"fmt"
	"strings"
	"time"
)

func LogTable(header []string, data map[string]time.Time) {
	// Find the longest key
	longestKey := 0
	for key := range data {
		if len(key) > longestKey {
			longestKey = len(key)
		}
	}

	// Calculate column widths
	col1Width := longestKey + 2
	col2Width := len(header[1]) + 2

	fmt.Println()
	fmt.Println(strings.Repeat("-", col1Width+col2Width+3))
	// Print header
	fmt.Printf("%-*s| %-*s\n", col1Width, header[0], col2Width, header[1])

	// Print separator
	fmt.Println(strings.Repeat("-", col1Width+col2Width+3))

	// Print data rows
	for key, value := range data {
		fmt.Printf("%-*s| %-*v\n", col1Width, key, col2Width, value)
	}
	fmt.Println(strings.Repeat("-", col1Width+col2Width+3))
}
