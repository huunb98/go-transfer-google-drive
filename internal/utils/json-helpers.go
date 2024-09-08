package utils

import (
	"encoding/json"
	"log"
)

// ParseJSON parses JSON data into the provided interface.
func ParseJSON(data string, v interface{}) error {
	return json.Unmarshal([]byte(data), v)
}

// MarshalJSON converts a Go value to a JSON-encoded string.
func MarshalJSON(v interface{}) (string, error) {
	jsonData, err := json.Marshal(v)
	if err != nil {
		return "", err
	}
	return string(jsonData), nil
}

// LogPrettyJSON logs a JSON representation of the given value with indentation.
func LogPrettyJSON(v interface{}) {
	jsonData, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		log.Printf("Error marshalling JSON: %v", err)
		return
	}
	log.Printf("JSON Output:\n%s", jsonData)
}
