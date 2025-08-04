package utils

import (
	"os"
	"testing"
)

func TestGenerateToken(t *testing.T) {
	// Set up test environment
	os.Setenv("TOKEN_SECRET", "test_secret")
	defer os.Unsetenv("TOKEN_SECRET")

	tests := []struct {
		name    string
		phone   string
		userId  int64
		wantErr bool
	}{
		{
			name:    "valid token generation",
			phone:   "1234567890",
			userId:  1,
			wantErr: false,
		},
		{
			name:    "empty phone",
			phone:   "",
			userId:  1,
			wantErr: false,
		},
		{
			name:    "zero user id",
			phone:   "1234567890",
			userId:  0,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := GernerateToken(tt.phone, tt.userId)
			if (err != nil) != tt.wantErr {
				t.Errorf("GernerateToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && token == "" {
				t.Error("GernerateToken() returned empty token")
			}
		})
	}
}

func TestVerifyToken(t *testing.T) {
	// Set up test environment
	os.Setenv("TOKEN_SECRET", "test_secret")
	defer os.Unsetenv("TOKEN_SECRET")

	// Generate a valid token for testing
	phone := "1234567890"
	userId := int64(1)
	token, err := GernerateToken(phone, userId)
	if err != nil {
		t.Fatalf("Failed to generate test token: %v", err)
	}

	tests := []struct {
		name    string
		token   string
		want    int64
		wantErr bool
	}{
		{
			name:    "valid token",
			token:   token,
			want:    userId,
			wantErr: false,
		},
		{
			name:    "invalid token",
			token:   "invalid.token.here",
			want:    0,
			wantErr: true,
		},
		{
			name:    "empty token",
			token:   "",
			want:    0,
			wantErr: true,
		},
		{
			name:    "malformed token",
			token:   "not.a.jwt.token",
			want:    0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := VerifyToken(tt.token)
			if (err != nil) != tt.wantErr {
				t.Errorf("VerifyToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("VerifyToken() = %v, want %v", got, tt.want)
			}
		})
	}
}
