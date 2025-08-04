package utils

import (
	"testing"
)

func TestHashPassword(t *testing.T) {
	tests := []struct {
		name     string
		password string
		wantErr  bool
	}{
		{
			name:     "valid password",
			password: "testPassword123",
			wantErr:  false,
		},
		{
			name:     "empty password",
			password: "",
			wantErr:  false,
		},
		{
			name:     "long password",
			password: "veryLongPasswordThatShouldStillWork123456789012345678901234567890",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := HashPassword(tt.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("HashPassword() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == "" {
				t.Error("HashPassword() returned empty string for valid password")
			}
		})
	}
}

func TestCheckPasswordHash(t *testing.T) {
	// First, create a hashed password to test with
	password := "testPassword123"
	hashedPassword, err := HashPassword(password)
	if err != nil {
		t.Fatalf("Failed to create test hash: %v", err)
	}

	tests := []struct {
		name           string
		password       string
		hashedPassword string
		want           bool
	}{
		{
			name:           "correct password",
			password:       password,
			hashedPassword: hashedPassword,
			want:           true,
		},
		{
			name:           "incorrect password",
			password:       "wrongPassword",
			hashedPassword: hashedPassword,
			want:           false,
		},
		{
			name:           "empty password",
			password:       "",
			hashedPassword: hashedPassword,
			want:           false,
		},
		{
			name:           "empty hash",
			password:       password,
			hashedPassword: "",
			want:           false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CheckPasswordHash(tt.password, tt.hashedPassword); got != tt.want {
				t.Errorf("CheckPasswordHash() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHashPasswordAndCheckPasswordHash(t *testing.T) {
	// Test that a hashed password can be verified
	password := "testPassword123"
	hashedPassword, err := HashPassword(password)
	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}

	if !CheckPasswordHash(password, hashedPassword) {
		t.Error("CheckPasswordHash() failed to verify a correctly hashed password")
	}
}
