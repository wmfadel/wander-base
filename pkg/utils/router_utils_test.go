package utils

import (
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/wmfadel/wander-base/internal/models"
)

func TestGetUserFromContext(t *testing.T) {
	// Set up test router
	gin.SetMode(gin.TestMode)

	// Test cases
	tests := []struct {
		name      string
		setupUser bool
		wantErr   bool
		setupFunc func(*gin.Context)
	}{
		{
			name:      "user exists in context",
			setupUser: true,
			wantErr:   false,
			setupFunc: func(c *gin.Context) {
				user := &models.User{ID: 1, Phone: "1234567890"}
				c.Set("user", *user)
			},
		},
		{
			name:      "no user in context",
			setupUser: false,
			wantErr:   true,
			setupFunc: func(c *gin.Context) {
				// Do nothing, context will be empty
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a test request
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			// Set up the context
			tt.setupFunc(c)

			// Test the function
			user, err := GetUserFromContext(c)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetUserFromContext() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if user == nil {
					t.Error("GetUserFromContext() returned nil user when error was not expected")
				}
			}
		})
	}
}

func TestGetUserIdFromContext(t *testing.T) {
	// Set up test router
	gin.SetMode(gin.TestMode)

	// Test cases
	tests := []struct {
		name      string
		setupUser bool
		wantErr   bool
		setupFunc func(*gin.Context)
	}{
		{
			name:      "user id exists in context",
			setupUser: true,
			wantErr:   false,
			setupFunc: func(c *gin.Context) {
				c.Set("userId", int64(1))
			},
		},
		{
			name:      "no user id in context",
			setupUser: false,
			wantErr:   true,
			setupFunc: func(c *gin.Context) {
				// Do nothing, context will be empty
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a test request
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			// Set up the context
			tt.setupFunc(c)

			// Test the function
			userId, err := GetUserIdFromContext(c)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetUserIdFromContext() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && userId == 0 {
				t.Error("GetUserIdFromContext() returned 0 when error was not expected")
			}
		})
	}
}
