package utils

import "time"

// GetCurrentTime returns the current time formatted as a string.
func GetCurrentTime() string {
	return time.Now().Format(time.RFC3339)
}

// ParseTime parses a time string into a time.Time object.
func ParseTime(timeStr string, layout string) (time.Time, error) {
	return time.Parse(layout, timeStr)
}
