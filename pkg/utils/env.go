package utils

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

func LoadEnv() error {
	if err := godotenv.Load("../.env"); err != nil {
		return fmt.Errorf("failed toload environment variables %w", err)
	}
	return nil
}

func GetFromEnv(key string) string {
	return os.Getenv(key)
}
