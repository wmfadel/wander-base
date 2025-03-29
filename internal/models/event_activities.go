package models

type EventActivities struct {
	EventID    int64 `gorm:"primaryKey;autoIncrement:false" json:"event_id"`
	ActivityID int64 `gorm:"primaryKey;autoIncrement:false" json:"activity_id"`
}

type RemoveActivitiesRequest struct {
	ActivityIDs []int64 `json:"activity_ids" binding:"required"`
}
