package main

import (
	"log"
	"os"
	"strconv"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/joho/godotenv"
	"github.com/snirkop89/lenslocked/models"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	host := os.Getenv("SMTP_HOST")
	portStr := os.Getenv("SMTP_PORT")
	username := os.Getenv("SMTP_USERNAME")
	password := os.Getenv("SMTP_PASSWORD")

	port, err := strconv.Atoi(portStr)
	if err != nil {
		log.Fatal(err)
	}

	config := models.SMTPConfig{
		Host:     host,
		Port:     port,
		Username: username,
		Password: password,
	}

	es := models.NewEmailService(config)
	err = es.ForgotPassword("me@example.com", "https://thisdomain.com/token?abc123")
	if err != nil {
		log.Fatal(err)
	}
}
