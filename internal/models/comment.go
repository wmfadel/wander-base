package models

import "time"

type Comment struct {
	ID        int64   `gorm:"primaryKey" json:"comment_id"`
	EventID   int64   `gorm:"index" json:"event_id"`
	UserID    int64   `gorm:"index" json:"user_id"`
	Content   string  `gorm:"not null" json:"content"`
	Score     float32 `gorm:"not null" json:"score"`
	Visible   bool    `gorm:"not null" json:"visible"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type CreateCommentRequest struct {
	Content string `json:"content" binding:"required"`
}
