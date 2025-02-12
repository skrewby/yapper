package utils

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Environment struct {
	ServerPort          string
	JWTSecret           string
	DatabaseCredentials DatabaseCredentials
}

type DatabaseCredentials struct {
	User string
	Pass string
	Host string
	Name string
	Port string
}

func GetEnvironmentVariables() Environment {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	serverPort := os.Getenv("SERVER_PORT")
	if serverPort == "" {
		log.Fatal("SERVER_PORT must have a value")
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET must have a value")
	}

	db := DatabaseCredentials{
		User: os.Getenv("DB_USER"),
		Pass: os.Getenv("DB_PASS"),
		Host: os.Getenv("DB_HOST"),
		Name: os.Getenv("DB_NAME"),
		Port: os.Getenv("DB_PORT"),
	}
	if db.User == "" || db.Pass == "" || db.Host == "" || db.Name == "" || db.Port == "" {
		log.Fatal("Database credentials not set in environment")
	}

	return Environment{
		ServerPort:          serverPort,
		JWTSecret:           jwtSecret,
		DatabaseCredentials: db,
	}
}
