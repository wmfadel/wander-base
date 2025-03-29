package models

type User struct {
	ID        int64  `gorm:"primaryKey" json:"id"`
	Phone     string `gorm:"not null;unique" json:"phone" binding:"required"`
	Password  string `gorm:"not null" json:"password" binding:"required,min=5"`
	FirstName string `gorm:"not null" json:"first_name" binding:"required"`
	LastName  string `gorm:"not null" json:"last_name" binding:"required"`
	Photo     string `json:"photo,omitempty"`
	Roles     []Role `gorm:"many2many:user_roles" json:"roles"`
}

func (u *User) Blocked() bool {
	return len(u.Roles) == 0
}
