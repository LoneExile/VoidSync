package utils

import (
	"os"
	"time"
)

func ConvertTimestamp(timestamp time.Time) time.Time {
	return timestamp.UTC().Truncate(time.Second)
}

func GetLocalTimestamp(path string) (time.Time, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return time.Time{}, err
	}
	return fileInfo.ModTime(), nil
}
