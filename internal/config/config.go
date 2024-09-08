package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

type Config struct {
	// OAuth2 configuration constants
	ClientID     string
	ClientSecret string
	RedirectURL  string

	// Server configuration constants
	ServerPort string

	// Database configuration constants
	MySqlHost     string
	MySqlPort     string
	MySqlUser     string
	MySqlPassword string
	MySqlDBName   string
}

func GetConfig() *Config {
	return &Config{
		ClientID:      os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret:  os.Getenv("GOOGLE_CLIENT_SECRET"),
		RedirectURL:   os.Getenv("GOOGLE_REDIRECT_URL"),
		ServerPort:    os.Getenv("SERVER_PORT"),
		MySqlHost:     os.Getenv("MYSQL_HOST"),
		MySqlPort:     os.Getenv("MYSQL_PORT"),
		MySqlUser:     os.Getenv("MYSQL_USERNAME"),
		MySqlPassword: os.Getenv("MYSQL_PASSWORD"),
		MySqlDBName:   os.Getenv("MYSQL_DATABASE_NAME"),
	}
}
