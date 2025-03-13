package repository

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/wmfadel/escape-be/internal/models"
)

type EventRepository struct {
	db *sql.DB
}

func NewEventRepository(db *sql.DB) *EventRepository {
	return &EventRepository{db: db}
}

func (repo *EventRepository) Save(event *models.Event) error {
	query := `INSERT INTO events(name, description, location, dateTime, user_id)
    VALUES ($1,$2,$3,$4,$5) RETURNING id`
	stmt, err := repo.db.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	var eventId int64
	err = stmt.QueryRow(event.Name, event.Description, event.Location, event.DateTime, event.UserID).Scan(&eventId)
	if err != nil {
		return fmt.Errorf("save event failed to insert: %w", err)
	}
	event.ID = eventId
	log.Printf("Created event %v", event)
	return nil
}

func (repo *EventRepository) Update(event *models.Event) error {
	query := `
	UPDATE events
	SET name = $1, description = $2, location = $3, dateTime = $4
	WHERE id = $5
	`
	stmt, err := repo.db.Prepare(query)
	if err != nil {
		return fmt.Errorf("Update event failed to prepare update event quer: %w", err)
	}
	defer stmt.Close()
	_, err = stmt.Exec(event.Name, event.Description, event.Location, event.DateTime, event.ID)
	if err != nil {
		return fmt.Errorf("executing update event query failed: %w", err)
	}
	return fmt.Errorf("executing update query failed %w", err)
}

func (repo *EventRepository) Delete(eventId int64) error {
	query := `DELETE FROM events WHERE id = $1`
	stmt, err := repo.db.Prepare(query)
	if err != nil {
		return fmt.Errorf("preparing event delete query failed %w", err)
	}
	defer stmt.Close()
	_, err = stmt.Exec(eventId)

	return fmt.Errorf("executing event delete query failed: %w", err)
}

func (repo *EventRepository) GetAllEvents() ([]models.Event, error) {
	query := "SELECT * FROM events"
	rows, err := repo.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query all events: %w", err)
	}
	defer rows.Close()

	var events = []models.Event{}

	for rows.Next() {
		var event models.Event
		err := rows.Scan(&event.ID, &event.Name, &event.Description, &event.Location, &event.DateTime, &event.UserID)

		if err != nil {
			return nil, fmt.Errorf("failed to scan event value: %w", err)
		}
		events = append(events, event)
	}
	return events, nil
}

func (repo *EventRepository) GetEventById(eventId int64) (*models.Event, error) {
	query := "SELECT * FROM events WHERE id = $1"
	row := repo.db.QueryRow(query, eventId)

	var event = models.Event{}
	err := row.Scan(&event.ID, &event.Name, &event.Description, &event.Location, &event.DateTime, &event.UserID)

	if err != nil {
		return nil, fmt.Errorf("failed to scan event after ID query: %w", err)
	}
	// TODO check if event is null and return error
	return &event, nil
}

func (repo *EventRepository) UpdatePartially(eventId int64, patch models.PatchEvent) error {

	if patch.IsEmpty() {
		return fmt.Errorf("no fields provided for update")
	}

	query, values, err := repo.buildUpdateQuery(eventId, patch)
	if err != nil {
		return err // Already wrapped with context
	}

	return repo.executeUpdateQuery(query, values)
}

// BuildUpdateQuery constructs the SQL query and values for the partial update
func (repo *EventRepository) buildUpdateQuery(id int64, patch models.PatchEvent) (string, []interface{}, error) {
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
func (repo *EventRepository) executeUpdateQuery(query string, values []interface{}) error {
	stmt, err := repo.db.Prepare(query)
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
func (repo *EventRepository) Register(userId, eventId int64) error {
	query := "INSERT INTO registrations (event_id, user_id) VALUES ($1, $2)"
	stmt, err := repo.db.Prepare(query)
	if err != nil {
		return fmt.Errorf("failed to prepare query for event registration: %w", err)
	}
	defer stmt.Close()

	result, err := stmt.Exec(eventId, userId)
	if err != nil {
		return fmt.Errorf("failed to execute registration query: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("registration failed: no rows affected")
	}
	return nil
}

func (repo *EventRepository) CancelRegister(userId int64, eventId int64) error {
	query := "DELETE FROM registration WHERE event_id=$1 AND user_id=$2"

	stmt, err := repo.db.Prepare(query)

	if err != nil {
		return fmt.Errorf("failed to prepare query for event unregistration %w", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(eventId, userId)
	return fmt.Errorf("failed to prepare query for event unregistration %w", err)
}
