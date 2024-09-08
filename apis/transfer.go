package apis

import (
	"encoding/json"
	"log"
	"net/http"

	transfer "transfer-folder-owner/pkg/transfer-owner"
)

// TransferRequest represents the request payload for transferring ownership.
// @Description Represents the request payload for transferring ownership.
type TransferRequest struct {
	// Email of the current owner
	// Required: true
	OriginEmail string `json:"origin_email" example:"currnetowner@example.com"`

	// Email of the new owner
	// Required: true
	NewOwnerEmail string `json:"new_owner_email" example:"newowner@example.com"`

	// FolderID of the folder whose ownership is to be transferred
	// Required: true
	FolderID string `json:"folder_id" example:"1BxiMVs0XRA5nFMdKvBdBZjgmUUqptlbs74OgvE2upms"`
}

type TransferOwnershipResponse struct {
	Message      string `json:"message" example:"Ownership transfer initiated"`
	Error        string `json:"error,omitempty" example:"Internal Server Error"`
	Duration     string `json:"duration,omitempty"`
	TotalFiles   int    `json:"total_files,omitempty"`
	TotalSuccess int    `json:"total_success,omitempty"`
	TotalErrors  int    `json:"total_errors,omitempty"`
}

// TransferHandler initiates the transfer of ownership of a Google Drive folder.
// @Summary Transfer ownership of a Google Drive folder
// @Description Initiates the transfer of ownership for a Google Drive folder to a specified email address.
// @Tags transfer
// @Accept  json
// @Produce  json
// @Param   request body TransferRequest true "Request payload for transferring ownership"
// @Success 200 {string} string "Ownership transfer initiated."
// @Failure 400 {string} string "Invalid request payload"
// @Failure 500 {string} string "Internal Server Error"
// @Router /transfer [post]
func TransferHandler(w http.ResponseWriter, r *http.Request) {
	var request TransferRequest

	// Parse JSON request bodyƯ
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	driveService, err := transfer.GetDriveService(request.OriginEmail)

	if err != nil {
		response := TransferOwnershipResponse{Message: "Error creating Drive service", Error: err.Error()}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Call the TransferOwnership function
	err = transfer.TransferOwnership(driveService, request.OriginEmail, request.FolderID, request.NewOwnerEmail)

	if err != nil {
		log.Printf("TransferOwnership error: %v", err)
		response := TransferOwnershipResponse{Message: "Failed to transfer ownership", Error: err.Error()}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	response := TransferOwnershipResponse{Message: "Ownership transfer completed successfully"}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

}
