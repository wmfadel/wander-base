package models

type User struct {
	ID       int64
	Phone    string `json:"phone" binding:"required"`
	Password string `json:"password" binding:"required,min=5"`
	Photo    string `json:"photo,omitempty"`
	Roles    []Role `json:"roles"`
}
