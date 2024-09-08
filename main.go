package main

import (
	"fmt"
	"log"
	"net/http"

	"transfer-folder-owner/apis"   // Import the apis package
	_ "transfer-folder-owner/docs" // Import the docs package for Swagger
	config "transfer-folder-owner/internal/config"
	"transfer-folder-owner/internal/database"

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title Example API
// @version 1.0
// @description This is a sample server using Gorilla Mux with multiple APIs.
// @host localhost:8080
// @BasePath /api/v1

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

func main() {
	log.Println("Starting the application")

	// load configuration
	var cfg = config.GetConfig()

	database.Connect()

	r := setupRoutes()

	fmt.Println("Server running on port", cfg.ServerPort)
	log.Fatal(http.ListenAndServe(":"+cfg.ServerPort, r))
}

func setupRoutes() *mux.Router {
	// Initialize the router
	r := mux.NewRouter()

	// Swagger documentation route
	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	// API routes
	api := r.PathPrefix("/api/v1").Subrouter()
	api.HandleFunc("/transfer", apis.TransferHandler).Methods("POST")
	api.HandleFunc("/oauth/google", apis.OAuthGoogleDrive).Methods("GET")
	api.HandleFunc("/oauth/google/callback", apis.OAuthGoogleDriveCallback).Methods("GET")

	return r
}
