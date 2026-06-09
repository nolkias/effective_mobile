package service

import (
	"context"
	"effective_mobile/internal/models"
	"effective_mobile/internal/repository"
	"fmt"
	"time"
)

type SubscriptionService struct {
	repo repository.SubscriptionRepo
}

func NewSubscriptionService(repo repository.SubscriptionRepo) *SubscriptionService {
	return &SubscriptionService{repo: repo}
}

func (s *SubscriptionService) Create(ctx context.Context, sub *models.Subscription) error {
	return s.repo.Create(ctx, sub)
}

func (s *SubscriptionService) Get(ctx context.Context, id string) (*models.Subscription, error) {
	return s.repo.Get(ctx, id)
}

func (s *SubscriptionService) GetList(ctx context.Context, limit, offset int) ([]*models.Subscription, int, error) {
	return s.repo.GetList(ctx, limit, offset)
}

func (s *SubscriptionService) Update(ctx context.Context, sub *models.Subscription) error {
	return s.repo.Update(ctx, sub)
}

func (s *SubscriptionService) Delete(ctx context.Context, id string) (bool, error) {
	return s.repo.Delete(ctx, id)
}

// TotalCost подсчитывает суммарную стоимость подписок за период.
// Для каждой подписки считается пересечение её периода с запрошенным диапазоном,
// результат умножается на цену за месяц.
func (s *SubscriptionService) TotalCost(ctx context.Context, startDate, endDate, userID, serviceName string) (int, error) {
	periodStart, err := parseMonthYear(startDate)
	if err != nil {
		return 0, fmt.Errorf("неверный формат start_date (ожидается MM-YYYY): %w", err)
	}
	periodEnd, err := parseMonthYear(endDate)
	if err != nil {
		return 0, fmt.Errorf("неверный формат end_date (ожидается MM-YYYY): %w", err)
	}

	subs, err := s.repo.GetByPeriod(ctx, startDate, endDate, userID, serviceName)
	if err != nil {
		return 0, err
	}

	var total int
	for _, sub := range subs {
		subStart, err := parseMonthYear(sub.StartDate)
		if err != nil {
			continue
		}

		// Конец подписки: либо указанная дата, либо конец запрошенного периода
		subEnd := periodEnd
		if sub.EndDate != nil && *sub.EndDate != "" {
			parsed, err := parseMonthYear(*sub.EndDate)
			if err == nil && parsed.Before(periodEnd) {
				subEnd = parsed
			}
		}

		// Начало пересечения
		overlapStart := periodStart
		if subStart.After(periodStart) {
			overlapStart = subStart
		}

		// Конец пересечения
		overlapEnd := subEnd

		months := monthsBetween(overlapStart, overlapEnd)
		if months > 0 {
			total += sub.Price * months
		}
	}

	return total, nil
}

// parseMonthYear разбирает строку формата MM-YYYY в начало месяца.
func parseMonthYear(s string) (time.Time, error) {
	return time.Parse("01-2006", s)
}

// monthsBetween считает количество полных месяцев включительно от start до end.
func monthsBetween(start, end time.Time) int {
	if end.Before(start) {
		return 0
	}
	months := (end.Year()-start.Year())*12 + int(end.Month()) - int(start.Month()) + 1
	return months
}
