package utils

import "strings"

// Capitalize capitalizes the first letter of a string.
func Capitalize(s string) string {
	if len(s) == 0 {
		return s
	}
	return strings.ToUpper(string(s[0])) + s[1:]
}

// TrimAndLower trims spaces and converts a string to lowercase.
func TrimAndLower(s string) string {
	return strings.ToLower(strings.TrimSpace(s))
}
