package models

import (
	"time"
)

type EventDestination struct {
	EventID       int64     `gorm:"primaryKey;autoIncrement:false" json:"event_id"`
	DestinationID int64     `gorm:"primaryKey;autoIncrement:false" json:"destination_id"`
	DateTime      time.Time `gorm:"not null" json:"datetime"`
}

type EventDestinationRequest struct {
	DestinationID int64     `json:"destination_id" binding:"required"`
	DateTime      time.Time `json:"datetime" binding:"required"`
}

type RemoveDestinationsRequest struct {
	DestinationIDs []int64 `json:"destination_ids" binding:"required"`
}
