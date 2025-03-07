package utils

import (
	"os"

	"github.com/joho/godotenv"
)

func LoadEnv() error {
	if err := godotenv.Load(); err != nil {
		return err
	}
	return nil
}

func GetFromEnv(key string) string {
	return os.Getenv(key)
}
