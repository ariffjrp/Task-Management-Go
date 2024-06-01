package configs

import (
	"log"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var GoogleOAuth2Config *oauth2.Config

func InitGoogleOAuth2Config() {
	clientID := os.Getenv("GOOGLE_CLIENT_ID")
	clientSecret := os.Getenv("GOOGLE_CLIENT_SECRET")

	if clientID == "" || clientSecret == "" {
		log.Fatal("Environment variables GOOGLE_CLIENT_ID and GOOGLE_CLIENT_SECRET must be set")
	}

	GoogleOAuth2Config = &oauth2.Config{
		RedirectURL:  "http://localhost:8080/v1/api/auth/login/oauth2/code/google",
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
		Endpoint:     google.Endpoint,
	}
}
