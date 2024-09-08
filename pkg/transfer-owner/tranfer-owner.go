package transfer_owner

import (
	"context"
	"fmt"
	"log"

	"transfer-folder-owner/internal/database"
	"transfer-folder-owner/internal/utils"
	oauth "transfer-folder-owner/pkg/oauth-google"

	"golang.org/x/oauth2"
	"google.golang.org/api/drive/v3"
)

type TransferRequest struct {
	Email    string `json:"email"`
	FolderID string `json:"folder_id"`
}

func TransferOwnership(driveService *drive.Service, originEmail, folderID, newOwnerEmail string) error {

	// Check if OriginEmail has access to the folder
	if !checkFolderPermission(driveService, folderID, originEmail) {
		return fmt.Errorf("OriginEmail %s does not have permission to access the folder", originEmail)
	}

	//  List files/folders and check ownership
	err := transferFilesIfOwner(driveService, folderID, originEmail, newOwnerEmail)
	if err != nil {
		return fmt.Errorf("failed to transfer ownership: %w", err)
	}

	return nil
}

// getFilesRecursively retrieves all files from a folder, including subfolders, with owner emails
func getFilesRecursively(driveService *drive.Service, folderID string) ([]*drive.File, error) {
	var allFiles []*drive.File

	// List all files and folders in the current folder
	fileList, err := driveService.Files.List().
		Q(fmt.Sprintf("'%s' in parents and trashed=false", folderID)).
		Fields("files(id, name, mimeType, owners, parents)").Do()

	if err != nil {
		return nil, fmt.Errorf("error listing files in folder %s: %v", folderID, err)
	}

	// Iterate through the retrieved files and folders
	for _, file := range fileList.Files {
		// Append file to result
		allFiles = append(allFiles, file)

		// Print file details including owner emails
		for _, owner := range file.Owners {
			log.Printf("File: %s (ID: %s) Owner: %s", file.Name, file.Id, owner.EmailAddress)
		}

		// If the file is a folder, recursively get files from the subfolder
		if file.MimeType == "application/vnd.google-apps.folder" {
			log.Printf("Recursing into folder: %s (ID: %s)", file.Name, file.Id)
			subFiles, err := getFilesRecursively(driveService, file.Id)
			if err != nil {
				return nil, err
			}
			allFiles = append(allFiles, subFiles...)
		}
	}

	return allFiles, nil
}

func transferFilesIfOwner(driveService *drive.Service, folderID, originEmail, newOwnerEmail string) error {
	// List files in the folder
	fileList, err := getFilesRecursively(driveService, folderID)
	if err != nil {
		return fmt.Errorf("failed to list files: %w", err)
	}

	log.Printf("File List: %v", fileList)

	for _, file := range fileList {
		log.Printf("File ID: %s, File Name: %s", file.Id, file.Name)
		// Check if OriginEmail is the owner of the file

		if isOwner(file, newOwnerEmail) {
			log.Println("File already owned by new owner")
			continue
		}

		if !isOwner(file, originEmail) {
			log.Printf("File %s is not owned by %s", file.Name, originEmail)
			continue
		}

		// Transfer ownership
		err := transferOwnership(driveService, file.Id, newOwnerEmail)
		if err != nil {
			log.Printf("Failed to transfer ownership of file: %s, Error: %v", file.Name, err)
		} else {
			log.Printf("Transferred ownership of file: %s", file.Name)
		}

	}

	return nil
}

// Helper to check if OriginEmail is the owner of the file
func isOwner(file *drive.File, originEmail string) bool {
	for _, owner := range file.Owners {
		if owner.EmailAddress == originEmail {
			return true
		}
	}
	return false
}

func checkFolderPermission(driveService *drive.Service, folderID, originEmail string) bool {
	folderName, ownerEmail, editEmails, err := getFolderDetails(driveService, folderID)
	if err != nil {
		log.Printf("Failed to get folder details: %v", err)
		return false
	}
	fmt.Printf("Folder Name: %s\n", folderName)
	fmt.Printf("Owner Email: %s\n", ownerEmail)
	fmt.Printf("Edit Permissions Emails: %v\n", editEmails)

	log.Println(err)
	if utils.ContainsString(editEmails, originEmail) {
		return true
	}

	return false
}

// getFolderDetails retrieves the folder's name and owner email from the Drive API
func getFolderDetails(driveService *drive.Service, folderID string) (string, string, []string, error) {
	// Retrieve folder metadata
	folder, err := driveService.Files.Get(folderID).Fields("name, owners").Do()
	// Get folder name and owner email
	folderName := folder.Name
	var ownerEmail string
	if len(folder.Owners) > 0 {
		ownerEmail = folder.Owners[0].EmailAddress
	} else {
		return folderName, "", nil, fmt.Errorf("no owner found for the folder")
	}

	// List permissions for the folder
	permissions, err := driveService.Permissions.List(folderID).Fields("permissions(emailAddress, role)").Do()
	if err != nil {
		return folderName, ownerEmail, nil, fmt.Errorf("error retrieving folder permissions: %w", err)
	}

	// Collect emails of users with edit (writer) permissions
	var editPermissionEmails []string
	for _, permission := range permissions.Permissions {
		if permission.Role == "writer" || permission.Role == "owner" {
			editPermissionEmails = append(editPermissionEmails, permission.EmailAddress)
		}
	}

	return folderName, ownerEmail, editPermissionEmails, nil
}

func HasAccess(originEmail string, folderID string) bool {
	// Use the Google Drive API to check if the OriginEmail has access to the folder
	driveService, err := GetDriveService(originEmail)
	if err != nil {
		return false
	}

	folder, err := driveService.Files.Get(folderID).Fields("owners").Do()
	if err != nil {
		return false
	}

	for _, owner := range folder.Owners {
		if owner.EmailAddress == originEmail {
			return true
		}
	}

	return false
}

func GetDriveService(originEmail string) (*drive.Service, error) {

	var config = oauth.Oauth2Config

	// Retrieve client secret from MySQL database
	var refreshToken string
	err := database.MySQL.QueryRow("SELECT refresh_token FROM oauth WHERE email = ?", originEmail).Scan(&refreshToken)
	if err != nil {
		return nil, err
	}

	// Create a token source from the refresh token
	token := &oauth2.Token{
		RefreshToken: refreshToken,
	}

	// Generate a new token source
	tokenSource := config.TokenSource(context.Background(), token)

	// Generate a new HTTP client with the token source
	client := oauth2.NewClient(context.Background(), tokenSource)

	// Create the Google Drive service
	driveService, err := drive.New(client)
	if err != nil {
		return nil, fmt.Errorf("failed to create drive service: %w", err)
	}

	return driveService, nil
}

// Transfer ownership of a file
func transferOwnership(driveService *drive.Service, fileID, newOwnerEmail string) error {
	permission := &drive.Permission{
		Type:         "user",
		Role:         "owner",
		EmailAddress: newOwnerEmail,
	}

	_, err := driveService.Permissions.Create(fileID, permission).TransferOwnership(true).Do()
	if err != nil {
		return fmt.Errorf("failed to transfer ownership: %w", err)
	}

	return nil
}
