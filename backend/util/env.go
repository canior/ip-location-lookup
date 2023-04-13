package util

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

func GetEnv(key string) string {
	env := os.Getenv("APP_ENV")

	file := "./.env"
	if env == "testing" {
		file = "../../.env.testing"
	}

	err := godotenv.Load(file)
	if err != nil {
		log.Fatalf("Error loading .env.%s file", env)
	}

	return os.Getenv(key)
}
