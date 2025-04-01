package service

import (
	"github.com/wmfadel/wander-base/internal/models"
	"github.com/wmfadel/wander-base/internal/repository"
)

type ActivityService struct {
	repo *repository.ActivityRepository
}

func NewActivityService(repo *repository.ActivityRepository) *ActivityService {
	return &ActivityService{repo: repo}
}

func (service *ActivityService) Save(activity *models.Activity) error {
	return service.repo.Save(activity)
}

func (service *ActivityService) GetAllActivities() ([]models.Activity, error) {
	return service.repo.GetAllActivities()
}

func (service *ActivityService) GetActivityById(id int64) (*models.Activity, error) {
	return service.repo.GetActivityById(id)
}

func (service *ActivityService) GetActivityBySlug(slug string) (*models.Activity, error) {
	return service.repo.GetActivityBySlug(slug)
}

func (service *ActivityService) Delete(activityId int64) error {
	return service.repo.Delete(activityId)
}
