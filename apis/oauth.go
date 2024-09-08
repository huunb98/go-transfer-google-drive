package apis

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	oauth "transfer-folder-owner/pkg/oauth-google"

	"golang.org/x/oauth2"

	database "transfer-folder-owner/internal/database"
	utils "transfer-folder-owner/internal/utils"
)

// UserInfo represents the user info response.
// @Description Represents the user info response.
type UserInfo struct {
	// Email of the user
	// Required: true
	Email string `json:"email"`

	// Picture of the user
	// Required: true
	Picture string `json:"picture"`

	// Expiry time of the token
	// Required: true
	Expiry int64 `json:"exp"`
}

// @Summary Redirect to Google OAuth2 authorization URL
// @Description Redirects the user to the Google OAuth2 authorization URL to obtain an access token
// @Tags oauth
// @Accept json
// @Produce json
// @Success 302 {string} string "Redirect URL"
// @Router /oauth/google [get]
func OAuthGoogleDrive(w http.ResponseWriter, r *http.Request) {
	url := oauth.Oauth2Config.AuthCodeURL("state", oauth2.AccessTypeOffline)
	log.Println(url)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"url": url})
}

// @Summary Google Drive OAuth2 callback
// @Description Exchanges the authorization code for an access token and retrieves user info
// @Tags oauth
// @Accept json
// @Produce json
// @Param code query string true "Authorization code"
// @Success 200 {object} UserInfo "User info"
// @Failure 400 {object} error "Missing code parameter"
// @Failure 500 {object} error "Failed to exchange token or retrieve user info"
// @Router /api/v1/oauth/google [get]
func OAuthGoogleDriveCallback(w http.ResponseWriter, r *http.Request) {
	// Get the authorization code from the request
	log.Println(r.URL)
	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "Missing code parameter", http.StatusBadRequest)
		return
	}

	token, err := oauth.Oauth2Config.Exchange(r.Context(), code)
	if err != nil {
		http.Error(w, "Failed to exchange token: "+err.Error(), http.StatusInternalServerError)
		return
	}

	utils.LogPrettyJSON(token)

	client := oauth.Oauth2Config.Client(r.Context(), token)
	response, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")

	if err != nil {
		http.Error(w, "Failed to get user info: "+err.Error(), http.StatusInternalServerError)
		return
	}

	utils.LogPrettyJSON(response.Body)

	var userInfo UserInfo

	if err := json.NewDecoder(response.Body).Decode(&userInfo); err != nil {
		http.Error(w, "Failed to decode user info: "+err.Error(), http.StatusInternalServerError)
		return
	}

	utils.LogPrettyJSON(userInfo)

	defer response.Body.Close()

	expiresAt := time.Time{}

	if userInfo.Expiry == 0 {
		expiresAt = time.Time{}
	} else {
		expiresAt = time.Unix(userInfo.Expiry, 0)
	}

	_, err = database.MySQL.Exec(`
        INSERT INTO oauth (avatar_url, provider, email, refresh_token, access_token, expires_at)
        VALUES (?, ?, ?, ?, ?, ?)
        ON DUPLICATE KEY UPDATE
        refresh_token = VALUES(refresh_token),
        access_token = VALUES(access_token),
        expires_at = VALUES(expires_at),
        updated_at = CURRENT_TIMESTAMP
    `, userInfo.Picture, "google", userInfo.Email, token.RefreshToken, token.AccessToken, expiresAt)

	if err != nil {
		http.Error(w, "Failed to save to database: "+err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "OAuth callback successful. Data saved.")
}
