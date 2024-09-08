package database

import (
	"database/sql"
	"fmt"
	"log"
	config "transfer-folder-owner/internal/config"

	_ "github.com/go-sql-driver/mysql" // Import the MySQL driver
)

// MySQL is a global variable that represents the database connection pool
var MySQL *sql.DB

// Connect creates a connection pool to the MySQL database
func Connect() {
	// Load the database configuration
	var cfg = config.GetConfig()

	// Open the connection pool
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s", cfg.MySqlUser, cfg.MySqlPassword, cfg.MySqlHost, cfg.MySqlDBName))
	if err != nil {
		log.Fatal(err)
	}

	// Set the maximum number of connections in the pool
	db.SetMaxOpenConns(10)

	// Set the maximum amount of time a connection can be idle
	db.SetConnMaxLifetime(0)

	// Set the maximum amount of time a connection can be open
	db.SetConnMaxIdleTime(0)

	// Ping the database to check the connection
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	// Assign the connection pool to the global variable
	MySQL = db
}
