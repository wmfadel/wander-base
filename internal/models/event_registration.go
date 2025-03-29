package models

// RegistrationStatus defines the possible status values
type RegistrationStatus string

// Define constants for the allowed status values
const (
	Registered          RegistrationStatus = "registered"
	PendingRegistration RegistrationStatus = "pending_registration"
	Cancelled           RegistrationStatus = "cancelled"
	PendingCancellation RegistrationStatus = "pending_cancellation"
)

// Registration model with composite key and status
type Registration struct {
	EventID int64              `gorm:"primaryKey;autoIncrement:false" json:"event_id"`
	UserID  int64              `gorm:"primaryKey;autoIncrement:false" json:"user_id"`
	Status  RegistrationStatus `gorm:"type:registration_status;not null;default:pending_registration" json:"status"`
	Event   Event              `gorm:"foreignKey:EventID;references:ID" json:"event,omitempty"`
	User    User               `gorm:"foreignKey:UserID;references:ID" json:"user,omitempty"`
}
