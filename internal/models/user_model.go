package models

type User struct {
	ID       int64
	Phone    string `binding:"required,phone"`
	Password string `binding:"required,min=5"`
	Roles    []Role `json:"roles"`
}
