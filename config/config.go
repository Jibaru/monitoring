package config

import (
	"errors"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	MongoURI  string
	JWTSecret string
	APIPort   string
	DBName    string
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

	return Config{
		MongoURI:  mongoURI,
		JWTSecret: jwtSecret,
		APIPort:   appPort,
		DBName:    "monitoringapp",
	}
}
