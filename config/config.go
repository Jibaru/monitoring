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
			log.Println(".env no encontrado, se utilizarán las variables de entorno por defecto")
		} else {
			log.Fatal("Error cargando .env", err)
		}
	}
	mongoURI, ok := os.LookupEnv("MONGODB_URI")
	if !ok {
		log.Fatal("MONGODB_URI no configurada")
	}

	jwtSecret, ok := os.LookupEnv("JWT_SECRET")
	if !ok {
		log.Fatal("JWT_SECRET no configurada")
	}

	appPort, ok := os.LookupEnv("PORT")
	if !ok {
		appPort = "8080"
	}

	mailAppPassword, ok := os.LookupEnv("MAIL_APP_PASSWORD")
	if !ok {
		log.Fatal("MAIL_APP_PASSWORD no configurada")
	}

	mailFromEmail, ok := os.LookupEnv("MAIL_FROM_EMAIL")
	if !ok {
		log.Fatal("MAIL_FROM_EMAIL no configurada")
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
		log.Fatal("GITHUB_CLIENT_ID no configurada")
	}

	githubClientSecret, ok := os.LookupEnv("GITHUB_CLIENT_SECRET")
	if !ok {
		log.Fatal("GITHUB_CLIENT_SECRET no configurada")
	}

	googleClientID, ok := os.LookupEnv("GOOGLE_CLIENT_ID")
	if !ok {
		log.Fatal("GOOGLE_CLIENT_ID no configurada")
	}

	googleClientSecret, ok := os.LookupEnv("GOOGLE_CLIENT_SECRET")
	if !ok {
		log.Fatal("GOOGLE_CLIENT_SECRET no configurada")
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
