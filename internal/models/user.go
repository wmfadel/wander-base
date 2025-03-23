package models

type User struct {
	ID        int64
	Phone     string `json:"phone" binding:"required"`
	Password  string `json:"password" binding:"required,min=5"`
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
	Photo     string `json:"photo,omitempty"`
	Roles     []Role `json:"roles"`
}
