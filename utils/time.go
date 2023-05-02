package utils

import "time"

func ConvertTimestamp(timestamp time.Time) time.Time {
	return timestamp.UTC().Truncate(time.Second)
}
