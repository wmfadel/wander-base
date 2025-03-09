package models

import (
	"log"
	"time"

	"github.com/wmfadel/escape-be/db"
)

type Event struct {
	ID          int64
	Name        string    `binding:"required"`
	Description string    `binding:"required"`
	Location    string    `binding:"required"`
	DateTime    time.Time `binding:"required"`
	UserID      int64
}

func (e *Event) Save() error {
	query := `INSERT INTO events(name, description, location, dateTime, user_id)
	VALUES ($1,$2,$3,$4,$5) RETURNING id`
	stmt, err := db.DB.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	var eventId int64
	err = stmt.QueryRow(e.Name, e.Description, e.Location, e.DateTime, e.UserID).Scan(&eventId)
	if err != nil {
		return err
	}
	e.ID = eventId
	log.Printf("Created event %v", e)
	return err
}

func (e Event) Update() error {
	query := `
	UPDATE events
	SET name = $1, description = $2, location = $3, dateTime = $4
	WHERE id = $5
	`
	stmt, err := db.DB.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()
	row := stmt.QueryRow(e.Name, e.Description, e.Location, e.DateTime, e.ID)
	return row.Err()
}

func (e Event) Delete() error {
	query := `DELETE FROM events WHERE id = $1`
	stmt, err := db.DB.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(e.ID)

	return err
}

func GetAllEvents() ([]Event, error) {
	query := "SELECT * FROM events"
	rows, err := db.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []Event

	for rows.Next() {
		var event Event
		err := rows.Scan(&event.ID, &event.Name, &event.Description, &event.Location, &event.DateTime, &event.UserID)

		if err != nil {
			return nil, err
		}
		events = append(events, event)
	}
	return events, nil
}

func GetEventById(id int64) (*Event, error) {
	query := "SELECT * FROM events WHERE id = $1"
	row := db.DB.QueryRow(query, id)

	var event Event
	err := row.Scan(&event.ID, &event.Name, &event.Description, &event.Location, &event.DateTime, &event.UserID)

	if err != nil {
		return nil, err
	}
	return &event, nil
}

func (e *Event) Register(userId int64) error {
	query := "INSERT INTO registration(event_id,user_id) VALUES ($1, $2)"

	stmt, err := db.DB.Prepare(query)

	if err != nil {
		return nil
	}
	defer stmt.Close()

	_, err = stmt.Exec(e.ID, userId)
	return err
}

func (e *Event) CancelRegister(userId int64) error {
	query := "DELETE FROM registration WHERE event_id=$1 AND user_id=$2"

	stmt, err := db.DB.Prepare(query)

	if err != nil {
		return nil
	}
	defer stmt.Close()

	_, err = stmt.Exec(e.ID, userId)
	return err
}
