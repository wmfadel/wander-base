package models

import (
	"fmt"
	"log"
	"strings"
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
		return fmt.Errorf("Save event failed to insert user %w", err)
	}
	e.ID = eventId
	log.Printf("Created event %v", e)
	return nil
}

func (e Event) Update() error {
	query := `
	UPDATE events
	SET name = $1, description = $2, location = $3, dateTime = $4
	WHERE id = $5
	`
	stmt, err := db.DB.Prepare(query)
	if err != nil {
		return fmt.Errorf("Update event failed to prepare update event quer: %w", err)
	}
	defer stmt.Close()
	_, err = stmt.Exec(e.Name, e.Description, e.Location, e.DateTime, e.ID)
	if err != nil {
		return fmt.Errorf("executing update event query failed: %w", err)
	}
	return fmt.Errorf("executing update query failed %w", err)
}

func (e Event) Delete() error {
	query := `DELETE FROM events WHERE id = $1`
	stmt, err := db.DB.Prepare(query)
	if err != nil {
		return fmt.Errorf("preparing event delete query failed %w", err)
	}
	defer stmt.Close()
	_, err = stmt.Exec(e.ID)

	return fmt.Errorf("executing event delete query failed: %w", err)
}

func GetAllEvents() ([]Event, error) {
	query := "SELECT * FROM events"
	rows, err := db.DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query all events: %w", err)
	}
	defer rows.Close()

	var events []Event

	for rows.Next() {
		var event Event
		err := rows.Scan(&event.ID, &event.Name, &event.Description, &event.Location, &event.DateTime, &event.UserID)

		if err != nil {
			return nil, fmt.Errorf("failed to scan event value: %w", err)
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
		return nil, fmt.Errorf("failed to scan event after ID query: %w", err)
	}
	return &event, nil
}

func (e *Event) Register(userId int64) error {
	query := "INSERT INTO registration(event_id,user_id) VALUES ($1, $2)"

	stmt, err := db.DB.Prepare(query)

	if err != nil {
		return fmt.Errorf("failed to prepare query for event registration %w", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(e.ID, userId)
	return fmt.Errorf("failed to execute registration query %w", err)
}

func (e *Event) CancelRegister(userId int64) error {
	query := "DELETE FROM registration WHERE event_id=$1 AND user_id=$2"

	stmt, err := db.DB.Prepare(query)

	if err != nil {
		return fmt.Errorf("failed to prepare query for event unregistration %w", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(e.ID, userId)
	return fmt.Errorf("failed to prepare query for event unregistration %w", err)
}

func (e Event) UpdatePartially(patch PatchEvent) error {

	if patch.IsEmpty() {
		return fmt.Errorf("no fields provided for update")
	}

	query, values, err := buildUpdateQuery(e.ID, patch)
	if err != nil {
		return err // Already wrapped with context
	}

	return executeUpdateQuery(query, values)
}

// BuildUpdateQuery constructs the SQL query and values for the partial update
func buildUpdateQuery(id int64, patch PatchEvent) (string, []interface{}, error) {
	var setClauses []string
	var values []interface{}
	paramIndex := 1

	// Build SET clauses and values for non-nil fields
	if patch.Name != nil {
		setClauses = append(setClauses, fmt.Sprintf("name = $%d", paramIndex))
		values = append(values, *patch.Name)
		paramIndex++
	}
	if patch.Description != nil {
		setClauses = append(setClauses, fmt.Sprintf("description = $%d", paramIndex))
		values = append(values, *patch.Description)
		paramIndex++
	}
	if patch.Location != nil {
		setClauses = append(setClauses, fmt.Sprintf("location = $%d", paramIndex))
		values = append(values, *patch.Location)
		paramIndex++
	}
	if patch.DateTime != nil {
		setClauses = append(setClauses, fmt.Sprintf("dateTime = $%d", paramIndex))
		values = append(values, *patch.DateTime)
		paramIndex++
	}

	if len(setClauses) == 0 {
		return "", nil, fmt.Errorf("no fields provided for update")
	}

	// Construct the query
	query := "UPDATE events SET " + strings.Join(setClauses, ", ") + " WHERE id = $" + fmt.Sprintf("%d", paramIndex)
	values = append(values, id)

	return query, values, nil
}

// ExecuteUpdateQuery executes the provided SQL query with the given values
func executeUpdateQuery(query string, values []interface{}) error {
	stmt, err := db.DB.Prepare(query)
	if err != nil {
		return fmt.Errorf("failed to prepare update query: %w", err)
	}
	defer stmt.Close()

	result, err := stmt.Exec(values...)
	if err != nil {
		return fmt.Errorf("failed to execute update query: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("no event found")
	}

	return nil
}
