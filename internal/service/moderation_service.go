package service

import (
	"context"
	"fmt"
	"log"

	"github.com/sashabaranov/go-openai"
	"github.com/wmfadel/wander-base/internal/models"
)

type ModerationService struct {
	openaiClient *openai.Client
	threshold    float32 // Threshold for hiding comments (e.g., 0.7)
}

// NewAuditService creates a new instance of AuditService
func NewAuditService(openaiClient *openai.Client, threshold float32) *ModerationService {
	return &ModerationService{
		openaiClient: openaiClient,
		threshold:    threshold,
	}
}

func (s *ModerationService) AuditComment(comment *models.Comment) error {
	// Call OpenAI's moderation API
	resp, err := s.openaiClient.Moderations(context.Background(), openai.ModerationRequest{
		Input: comment.Content,
		Model: openai.ModerationTextStable,
	})
	log.Println(resp)
	if err != nil {
		return fmt.Errorf("failed to call OpenAI moderation API: %w", err)
	}

	if len(resp.Results) == 0 {
		return fmt.Errorf("no moderation results returned")
	}

	result := resp.Results[0]

	// List all category scores
	scores := []float32{
		result.CategoryScores.Hate,
		result.CategoryScores.HateThreatening,
		result.CategoryScores.SelfHarm,
		result.CategoryScores.Sexual,
		result.CategoryScores.SexualMinors,
		result.CategoryScores.Violence,
		result.CategoryScores.ViolenceGraphic,
	}

	// Find the maximum score
	maxScore := float32(0.0)
	for _, score := range scores {
		if score > maxScore {
			maxScore = score
		}
	}
	comment.Score = maxScore

	// Determine visibility based on threshold
	comment.Visible = true
	if result.CategoryScores.Hate > s.threshold ||
		result.CategoryScores.HateThreatening > s.threshold ||
		result.CategoryScores.SelfHarm > s.threshold ||
		result.CategoryScores.Sexual > s.threshold ||
		result.CategoryScores.SexualMinors > s.threshold ||
		result.CategoryScores.Violence > s.threshold ||
		result.CategoryScores.ViolenceGraphic > s.threshold {
		comment.Visible = false
	}

	return nil
}
