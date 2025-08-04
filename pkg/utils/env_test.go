package utils

import (
	"os"
	"testing"
)

func TestLoadEnv(t *testing.T) {
	// Test with existing .env file
	err := LoadEnv()
	if err != nil {
		t.Errorf("LoadEnv() error = %v", err)
	}

	// Test with non-existent .env file
	os.Rename(".env", ".env.backup")
	defer os.Rename(".env.backup", ".env")

	err = LoadEnv()
	if err == nil {
		t.Error("LoadEnv() expected error for missing .env file")
	}
}

func TestGetFromEnv(t *testing.T) {
	// Set up test environment variable
	testKey := "TEST_ENV_VAR"
	testValue := "test_value"
	os.Setenv(testKey, testValue)
	defer os.Unsetenv(testKey)

	// Test getting existing environment variable
	got := GetFromEnv(testKey)
	if got != testValue {
		t.Errorf("GetFromEnv() = %v, want %v", got, testValue)
	}

	// Test getting non-existent environment variable
	got = GetFromEnv("NON_EXISTENT_VAR")
	if got != "" {
		t.Errorf("GetFromEnv() = %v, want empty string", got)
	}
}
