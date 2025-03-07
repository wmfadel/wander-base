package models

import (
	"errors"

	"githuv.com/wmfadel/go_events/db"
	"githuv.com/wmfadel/go_events/utils"
)

type User struct {
	ID       int64
	Email    string `binding:"required"`
	Password string `binding:"required"`
}

func (u *User) Save() error {
	query := "INSERT INTO users(email, password) VALUES (?, ?)"

	stmt, err := db.DB.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	hashedPassword, err := utils.HashPassword(u.Password)
	if err != nil {
		return err
	}
	res, err := stmt.Exec(u.Email, hashedPassword)
	if err != nil {
		return err
	}
	userId, err := res.LastInsertId()
	if err != nil {
		return err
	}
	u.ID = userId
	return err
}

func (u *User) ValidateCredintials() error {
	query := "SELECT id, password FROM users WHERE email = ?"

	row := db.DB.QueryRow(query, u.Email)
	var storedPassword string

	err := row.Scan(&u.ID, &storedPassword)

	if err != nil {
		return errors.New("wrong credintials")
	}

	isValidPassword := utils.CheckPasswordHash(u.Password, storedPassword)

	if !isValidPassword {
		return errors.New("wrong credintials")
	}
	return nil
}
