package repository

import (
	"fmt"

	"github.com/wmfadel/wander-base/internal/models"
	"github.com/wmfadel/wander-base/internal/models/requests"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type EventRepository struct {
	db        *gorm.DB
	photoRepo *EventPhotoRepository
}

func NewEventRepository(db *gorm.DB, photoRepo *EventPhotoRepository) *EventRepository {
	return &EventRepository{db: db, photoRepo: photoRepo}
}

func (repo *EventRepository) Save(event *models.Event) error {
	if event == nil {
		return fmt.Errorf("event is nil")
	}
	if event.IsEmpty() {
		return fmt.Errorf("event is empty")
	}
	result := repo.db.Create(event)
	if result.Error != nil {
		return result.Error
	}

	repo.db.Last(&event)
	return nil
}

func (repo *EventRepository) SetDestinations(destinations []models.EventDestinationRequest, eventID int64) error {
	return repo.db.Transaction(func(tx *gorm.DB) error {
		// Prepare new event_destination records
		var eventDestinations []models.EventDestination
		for _, req := range destinations {
			eventDestinations = append(eventDestinations, models.EventDestination{
				EventID:       eventID,
				DestinationID: req.DestinationID,
				DateTime:      req.DateTime,
			})
		}

		// Insert new destinations, ignoring duplicates
		if len(eventDestinations) > 0 {
			err := tx.Clauses(clause.OnConflict{DoNothing: true}).Create(&eventDestinations).Error
			if err != nil {
				return fmt.Errorf("failed to append destinations to event %d: %w", eventID, err)
			}
		}

		return nil
	})
}

func (repo *EventRepository) RemoveDestinations(destinationIDs []int64, eventID int64) error {
	return repo.db.Transaction(func(tx *gorm.DB) error {
		// Remove specified destinations for the event
		result := tx.Where("event_id = ? AND destination_id IN ?", eventID, destinationIDs).
			Delete(&models.EventDestination{})
		if result.Error != nil {
			return fmt.Errorf("failed to remove destinations from event %d: %w", eventID, result.Error)
		}

		// Optional: Check if any rows were affected (not strictly necessary unless you want to fail on no-op)
		if result.RowsAffected == 0 {
			return fmt.Errorf("no destinations removed for event %d (none matched the provided IDs)", eventID)
		}

		return nil
	})
}

func (repo *EventRepository) AddActivities(activityIDs []int64, eventID int64) error {
	return repo.db.Transaction(func(tx *gorm.DB) error {
		// Prepare new event_activities records
		var eventActivities []models.EventActivities
		for _, activityID := range activityIDs {
			eventActivities = append(eventActivities, models.EventActivities{
				EventID:    eventID,
				ActivityID: activityID,
			})
		}

		// Insert new activities, ignoring duplicates
		if len(eventActivities) > 0 {
			err := tx.Clauses(clause.OnConflict{DoNothing: true}).Create(&eventActivities).Error
			if err != nil {
				return fmt.Errorf("failed to add activities to event %d: %w", eventID, err)
			}
		}

		return nil
	})
}

func (repo *EventRepository) RemoveActivities(activityIDs []int64, eventID int64) error {
	return repo.db.Transaction(func(tx *gorm.DB) error {
		// Remove specified activities for the event
		result := tx.Where("event_id = ? AND activity_id IN ?", eventID, activityIDs).
			Delete(&models.EventActivities{})
		if result.Error != nil {
			return fmt.Errorf("failed to remove activities from event %d: %w", eventID, result.Error)
		}

		// Optional: Check if any rows were affected
		if result.RowsAffected == 0 {
			return fmt.Errorf("no activities removed for event %d (none matched the provided IDs)", eventID)
		}

		return nil
	})
}

func (repo *EventRepository) Update(event *models.Event) error {

	result := repo.db.Save(event)
	if result.Error != nil {
		return fmt.Errorf("executing update event query failed: %w", result.Error)
	}
	return nil
}

func (repo *EventRepository) Delete(eventId int64) error {
	result := repo.db.Delete(&models.Event{}, eventId)

	if result.Error != nil {
		return fmt.Errorf("failed to delete event %d: %w", eventId, result.Error)
	}

	return nil
}

func (repo *EventRepository) GetAllEvents() ([]models.Event, error) {
	var events = []models.Event{}

	result := repo.db.Find(&events)

	if result.Error != nil {
		return nil, fmt.Errorf("failed to get all events: %w", result.Error)
	}
	return events, nil
}

func (repo *EventRepository) GetEventById(id int64) (*models.Event, error) {
	var event models.Event

	// Fetch event with associations
	err := repo.db.
		Preload("User").             // Load associated User
		Preload("Destinations").     // Load associated Destinations via event_destinations
		Preload("Activities").       // Load associated Activities via event_activities
		Preload("Photos").           // Load associated Photos
		First(&event, "id = ?", id). // Fetch event by ID
		Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil // Return nil, nil if event not found
		}
		return nil, fmt.Errorf("failed to get event %d: %w", id, err)
	}

	return &event, nil
}

func (repo *EventRepository) UpdatePartially(eventID int64, patch requests.PatchEvent) error {
	if patch.IsEmpty() {
		return fmt.Errorf("no fields provided for update")
	}

	updates := make(map[string]interface{})
	if patch.Name != nil {
		updates["name"] = *patch.Name
	}
	if patch.Description != nil {
		updates["description"] = *patch.Description
	}
	if patch.Location != nil {
		updates["location"] = *patch.Location
	}
	if patch.DateTime != nil {
		updates["date_time"] = *patch.DateTime
	}

	result := repo.db.Model(&models.Event{}).Where("id = ?", eventID).Updates(updates)
	if result.Error != nil {
		return fmt.Errorf("failed to update event %d: %w", eventID, result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("no event found with id %d", eventID)
	}

	return nil
}
