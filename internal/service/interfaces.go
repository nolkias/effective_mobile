package service

import (
	"context"
	"effective_mobile/internal/models"
)

type SubscriptionSvc interface {
	Create(ctx context.Context, sub *models.Subscription) error
	Get(ctx context.Context, id string) (*models.Subscription, error)
	GetList(ctx context.Context, limit, offset int) ([]*models.Subscription, int, error)
	Update(ctx context.Context, sub *models.Subscription) error
	Delete(ctx context.Context, id string) (bool, error)
	TotalCost(ctx context.Context, startDate, endDate, userID, serviceName string) (int, error)
}
