package models

import (
	"time"
)

type EventPhoto struct {
	ID         int64     `gorm:"primaryKey" json:"_"`
	EventID    int64     `gorm:"index" json:"__"`
	URL        string    `gorm:"not null" json:"url"`
	UploadedAt time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"___"`
	Event      Event     `gorm:"foreignKey:EventID;references:ID"`
}
