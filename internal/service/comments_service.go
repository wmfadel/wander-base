package service

import (
	"github.com/wmfadel/wander-base/internal/models"
	"github.com/wmfadel/wander-base/internal/repository"
)

type CommentService struct {
	repo         *repository.CommentRepository
	auditService *ModerationService
}

func NewCommentService(repo *repository.CommentRepository, auditService *ModerationService) *CommentService {
	return &CommentService{
		repo:         repo,
		auditService: auditService,
	}
}

func (service *CommentService) Create(comment *models.Comment) error {
	// TODO: Uncomment this when we have a moderation service
	// if err := service.auditService.AuditComment(comment); err != nil {
	// 	return fmt.Errorf("failed to audit comment: %w", err)
	// }
	return service.repo.Create(comment)
}

func (service *CommentService) Delete(commentId int64) error {
	return service.repo.Delete(commentId)
}

func (service *CommentService) GetCommentById(commentId int64) (*models.Comment, error) {
	return service.repo.GetCommentById(commentId)
}

func (service *CommentService) GetEventComments(EventID int64) ([]models.Comment, error) {
	return service.repo.GetEventComments(EventID)
}
