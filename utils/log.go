package utils

import (
	"fmt"
	"strings"
	"voidsync/storage"
)

func LogTable(header []string, data map[string]storage.FileInfo) {
	longestFirstKey := 0
	for key := range data {
		if len(key) > longestFirstKey {
			longestFirstKey = len(key)
		}
	}
	longestSecondKey := 0
	for _, value := range data {
		if len(value.ETag) > longestSecondKey {
			longestSecondKey = len(value.ETag)
		}
	}

	col1Width := longestFirstKey + 2
	col2Width := longestSecondKey + 2
	col3Width := len(header[2]) + 2

	fmt.Println()
	fmt.Println(strings.Repeat("-", col1Width+col2Width+col3Width+20))
	fmt.Printf("📂 %-*s| %-*s| %-*s\n", col1Width, header[0], col2Width, header[1], col3Width, header[2])

	fmt.Println(strings.Repeat("-", col1Width+col2Width+col3Width+20))

	for key, value := range data {
		fmt.Printf("📃 %-*s| %-*v| %-*v\n", col1Width, key, col2Width, value.ETag, col3Width, value.Timestamp)
	}
	fmt.Println(strings.Repeat("-", col1Width+col2Width+col3Width+20))
}
