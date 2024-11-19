package database

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func LoadEnv() string {
	err := godotenv.Load("D:\\Amir\\Code\\System_Group\\Accounting_system_proj\\configs\\.env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	dbHost := os.Getenv("DB_HOST")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbPort := os.Getenv("DB_PORT")
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", dbHost, dbUser, dbPassword, dbName, dbPort)
}
