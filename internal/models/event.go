package models

import (
	"time"
)

// / Events are created without photos, destinations or activities, they are added later on.
type Event struct {
	ID           int64         `gorm:"primaryKey" json:"id"`
	Name         string        `gorm:"not null" json:"name"`
	Description  string        `json:"description"`
	Location     string        `gorm:"not null" json:"location"`
	DateTime     time.Time     `gorm:"not null" json:"date_time"`
	UserID       int64         `gorm:"index;not null" json:"user_id"`
	Photos       []EventPhoto  `gorm:"foreignKey:EventID" json:"photos,omitempty"`
	Destinations []Destination `gorm:"many2many:event_destinations"`
	Activities   []Activity    `gorm:"many2many:event_activities"`
}

func (e Event) PhotosUrls() []string {
	photoUrls := []string{}
	for _, photo := range e.Photos {
		photoUrls = append(photoUrls, photo.URL)
	}
	return photoUrls
}

func (e Event) IsEmpty() bool {
	if e.Name == "" && e.Description == "" && e.Location == "" && e.DateTime.IsZero() {
		return true
	}
	return false
}
