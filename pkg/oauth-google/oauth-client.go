package oauth_google

import (
	config "transfer-folder-owner/internal/config"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var cfg = config.GetConfig()

var Oauth2Config = &oauth2.Config{
	ClientID:     cfg.ClientID,
	ClientSecret: cfg.ClientSecret,
	RedirectURL:  cfg.RedirectURL,
	Scopes: []string{"https://www.googleapis.com/auth/userinfo.profile",
		"https://www.googleapis.com/auth/userinfo.email",
		"https://www.googleapis.com/auth/drive",
		"https://www.googleapis.com/auth/drive.file",
		"https://www.googleapis.com/auth/drive.appfolder",
	},
	Endpoint: google.Endpoint,
}
