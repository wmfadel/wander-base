package utils

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/wmfadel/escape-be/internal/models"
)

func GetUserFromContext(context *gin.Context) (*models.User, error) {
	value, exists := context.Get("user")
	if !exists {
		return nil, fmt.Errorf("user not found in context")
	}

	// Type assert the value to models.Event
	user, ok := value.(models.User)
	if !ok {
		return nil, fmt.Errorf("invalid user data in context %v", user)
	}
	return &user, nil
}

func GetEventFromContext(context *gin.Context) (*models.Event, error) {
	value, exists := context.Get("event")
	if !exists {
		return nil, fmt.Errorf("event not found in context")
	}

	// Type assert the value to models.Event
	event, ok := value.(models.Event)
	if !ok {
		return nil, fmt.Errorf("invalid event data in context %v", event)
	}

	return &event, nil
}
