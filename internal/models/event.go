package models

import "time"

type Event struct {
	ID          int64     `json:"id" form:"id"`
	Name        string    `json:"name" form:"name" binding:"required"`
	Description string    `json:"description" form:"description"`
	Location    string    `json:"location" form:"location" binding:"required"`
	DateTime    time.Time `json:"dateTime" form:"dateTime" binding:"required"`
	UserID      int64     `json:"user_id" form:"user_id"`
	Photos      []string  `json:"photos,omitempty" form:"photos,omitempty"` // Populated after upload
}
