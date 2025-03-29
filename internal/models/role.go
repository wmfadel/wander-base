package models

type Role struct {
	ID          int64  `gorm:"primaryKey" json:"id"`
	Name        string `gorm:"not null;unique" json:"name"`
	Description string `gorm:"not null" json:"description"`
	Default     bool   `gorm:"column:default_role;not null" json:"default"`
}
