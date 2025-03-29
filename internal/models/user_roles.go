package models

type UserRole struct {
	UserID int64 `gorm:"primaryKey;autoIncrement:false" json:"user_id"`
	RoleID int64 `gorm:"primaryKey;autoIncrement:false" json:"role_id"`
	User   User  `gorm:"foreignKey:UserID;references:ID" json:"user,omitempty"`
	Role   Role  `gorm:"foreignKey:RoleID;references:ID" json:"role,omitempty"`
}
