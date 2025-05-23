package models

// Destination model with GORM tags
type Destination struct {
	ID          int64  `gorm:"primaryKey" json:"id"`
	Name        string `gorm:"not null" json:"name"`
	Location    string `gorm:"not null" json:"location"`
	Description string `json:"description,omitempty"`
}
