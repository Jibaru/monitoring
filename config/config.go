package config

import (
	"errors"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	APIBaseURI         string
	WebBaseURI         string
	MongoURI           string
	JWTSecret          string
	APIPort            string
	DBName             string
	MailAppPassword    string
	MailFromEmail      string
	GithubClientID     string
	GithubClientSecret string
	GoogleClientID     string
	GoogleClientSecret string
}

func Load() Config {
	if err := godotenv.Load(); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			log.Println(".env not found, using environment variables as default")
		} else {
			log.Fatal("error loading .env", err)
		}
	}
	mongoURI, ok := os.LookupEnv("MONGODB_URI")
	if !ok {
		log.Fatal("MONGODB_URI not configured")
	}

	jwtSecret, ok := os.LookupEnv("JWT_SECRET")
	if !ok {
		log.Fatal("JWT_SECRET not configured")
	}

	appPort, ok := os.LookupEnv("PORT")
	if !ok {
		appPort = "8080"
	}

	mailAppPassword, ok := os.LookupEnv("MAIL_APP_PASSWORD")
	if !ok {
		log.Fatal("MAIL_APP_PASSWORD not configured")
	}

	mailFromEmail, ok := os.LookupEnv("MAIL_FROM_EMAIL")
	if !ok {
		log.Fatal("MAIL_FROM_EMAIL not configured")
	}

	APIBaseURI, ok := os.LookupEnv("API_BASE_URI")
	if !ok {
		APIBaseURI = "http://localhost:8080"
	}

	webBaseURI, ok := os.LookupEnv("WEB_BASE_URI")
	if !ok {
		webBaseURI = "http://localhost:5173"
	}

	githubClientID, ok := os.LookupEnv("GITHUB_CLIENT_ID")
	if !ok {
		log.Fatal("GITHUB_CLIENT_ID not configured")
	}

	githubClientSecret, ok := os.LookupEnv("GITHUB_CLIENT_SECRET")
	if !ok {
		log.Fatal("GITHUB_CLIENT_SECRET not configured")
	}

	googleClientID, ok := os.LookupEnv("GOOGLE_CLIENT_ID")
	if !ok {
		log.Fatal("GOOGLE_CLIENT_ID not configured")
	}

	googleClientSecret, ok := os.LookupEnv("GOOGLE_CLIENT_SECRET")
	if !ok {
		log.Fatal("GOOGLE_CLIENT_SECRET not configured")
	}

	return Config{
		APIBaseURI:         APIBaseURI,
		WebBaseURI:         webBaseURI,
		MongoURI:           mongoURI,
		JWTSecret:          jwtSecret,
		APIPort:            appPort,
		DBName:             "monitoringapp",
		MailAppPassword:    mailAppPassword,
		MailFromEmail:      mailFromEmail,
		GithubClientID:     githubClientID,
		GithubClientSecret: githubClientSecret,
		GoogleClientID:     googleClientID,
		GoogleClientSecret: googleClientSecret,
	}
}
