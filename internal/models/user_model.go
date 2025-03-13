package models

type User struct {
	ID       int64
	Email    string `binding:"required,email"`
	Password string `binding:"required,min=5"`
}
