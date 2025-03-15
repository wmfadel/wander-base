package repository

import (
	"database/sql"
	"fmt"
	"mime/multipart"
	"strings"

	"github.com/lib/pq"
	"github.com/wmfadel/escape-be/pkg/utils"
)

type EventPhotoRepository struct {
	db      *sql.DB
	storage *utils.Storage
}

func NewEventPhotoRepository(db *sql.DB, storage *utils.Storage) *EventPhotoRepository {
	return &EventPhotoRepository{db: db, storage: storage}
}

func (repo *EventPhotoRepository) AddPhotos(eventID int64, photos []*multipart.FileHeader) error {
	query := "INSERT INTO event_photos (event_id, photo_url) VALUES ($1, $2)"
	stmt, err := repo.db.Prepare(query)
	if err != nil {
		return fmt.Errorf("failed to prepare query for adding event photos: %w", err)
	}
	defer stmt.Close()

	for _, photo := range photos {
		url, err := repo.storage.UploadFile(photo, "events", eventID)
		if err != nil {
			return fmt.Errorf("failed to upload photo: %w", err)
		}

		_, err = stmt.Exec(eventID, url)
		if err != nil {
			if pqErr, ok := err.(*pq.Error); ok {
				switch pqErr.Code {
				case "23503": // foreign_key_violation
					return fmt.Errorf("event %d does not exist: %w", eventID, err)
				default:
					return fmt.Errorf("database error (code %s): %w", pqErr.Code, err)
				}
			}
			return fmt.Errorf("failed to add photo to event %d: %w", eventID, err)
		}
	}
	return nil
}

func (repo *EventPhotoRepository) GetPhotos(eventID int64) ([]string, error) {
	query := "SELECT photo_url FROM event_photos WHERE event_id = $1"
	rows, err := repo.db.Query(query, eventID)
	if err != nil {
		return nil, fmt.Errorf("failed to query photos for event %d: %w", eventID, err)
	}
	defer rows.Close()

	var photos []string
	for rows.Next() {
		var url string
		if err := rows.Scan(&url); err != nil {
			return nil, fmt.Errorf("failed to scan photo URL: %w", err)
		}
		photos = append(photos, url)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating photo rows: %w", err)
	}
	return photos, nil
}

func (repo *EventPhotoRepository) DeletePhotos(eventID int64, urls []string) error {
	if len(urls) == 0 {
		return nil // Nothing to delete
	}

	// Prepare query with dynamic placeholders for URLs
	placeholders := make([]string, len(urls))
	args := make([]interface{}, len(urls)+1)
	args[0] = eventID
	for i, url := range urls {
		placeholders[i] = fmt.Sprintf("$%d", i+2) // $2, $3, etc.
		args[i+1] = url
	}
	query := fmt.Sprintf("DELETE FROM event_photos WHERE event_id = $1 AND photo_url IN (%s)",
		strings.Join(placeholders, ","))

	// Execute deletion from database
	result, err := repo.db.Exec(query, args...)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code {
			case "23503": // foreign_key_violation (unlikely here due to CASCADE)
				return fmt.Errorf("event %d does not exist: %w", eventID, err)
			default:
				return fmt.Errorf("database error (code %s): %w", pqErr.Code, err)
			}
		}
		return fmt.Errorf("failed to delete photos from event %d: %w", eventID, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}

	// Delete files from filesystem
	for _, url := range urls {
		if err := repo.storage.DeleteFile(url); err != nil {
			fmt.Printf("warning: failed to delete file %s: %v\n", url, err)
			// Continue despite failure—DB is already updated
		}
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no photos found for event %d matching provided URLs", eventID)
	}

	return nil
}
