package repository

import (
	"fmt"

	"github.com/wmfadel/wander-base/internal/models"
	"gorm.io/gorm"
)

type CommentRepository struct {
	db *gorm.DB
}

func NewCommentRepository(db *gorm.DB) *CommentRepository {
	return &CommentRepository{db: db}
}

func (repo *CommentRepository) Create(comment *models.Comment) error {
	result := repo.db.Create(comment)
	if result.Error != nil {
		return fmt.Errorf("failed to save comment: %w", result.Error)
	}
	return nil
}

func (repo *CommentRepository) Delete(commentId int64) error {
	result := repo.db.Delete(&models.Comment{}, commentId)
	if result.Error != nil {
		return fmt.Errorf("failed to delete comment %d: %w", commentId, result.Error)
	}
	return nil
}

func (repo *CommentRepository) GetCommentById(commentId int64) (*models.Comment, error) {
	comment := models.Comment{ID: commentId}
	result := repo.db.First(&comment)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get comment %d: %w", commentId, result.Error)
	}
	return &comment, nil
}

func (repo *CommentRepository) GetEventComments(EventID int64) ([]models.Comment, error) {
	var comments []models.Comment
	result := repo.db.Where("event_id = ?", EventID).Find(&comments)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to get comments for event %d: %w", EventID, result.Error)
	}
	return comments, nil
}
