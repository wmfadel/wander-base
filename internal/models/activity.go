package models

type Activity struct {
	ID          int64  `gorm:"primaryKey" json:"id"`
	Name        string `gorm:"not null" json:"name"`
	Slug        string `gorm:"not null" json:"slug"`
	Description string `json:"description,omitempty"`
}
