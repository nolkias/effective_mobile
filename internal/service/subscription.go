package service

import (
	"effective_mobile/internal/models"
	"effective_mobile/internal/repository"
)

type SubscriptionService struct {
	repo *repository.SubscriptionRepository
}

func NewSubscriptionService(repo *repository.SubscriptionRepository) *SubscriptionService {
	return &SubscriptionService{repo: repo}
}

func (s *SubscriptionService) Create(sub *models.Subscription) error {
	return s.repo.Create(sub)
}

func (s *SubscriptionService) Update(sub *models.Subscription) error {
	return s.repo.Update(sub)
}

func (s *SubscriptionService) Get(id string) (*models.Subscription, error) {
	return s.repo.Get(id)
}

func (s *SubscriptionService) GetList(limit, offset int) ([]*models.Subscription, int, error) {
	return s.repo.GetList(limit, offset)
}

func (s *SubscriptionService) Delete(id string) error {
	return s.repo.Delete(id)
}

func (s *SubscriptionService) TotalCost(startDate, endDate, userID, serviceName string) (int, error) {
	return s.repo.TotalCost(startDate, endDate, userID, serviceName)
}
