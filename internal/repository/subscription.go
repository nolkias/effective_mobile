package repository

import (
	"context"
	"database/sql"
	"effective_mobile/internal/models"
	"errors"
)

type SubscriptionRepository struct {
	db *sql.DB
}

func NewSubscriptionRepository(db *sql.DB) *SubscriptionRepository {
	return &SubscriptionRepository{db: db}
}

func (r *SubscriptionRepository) Create(ctx context.Context, sub *models.Subscription) error {
	query := `INSERT INTO subscriptions (id, service_name, price, user_id, start_date, end_date, created_at, updated_at)
              VALUES ($1, $2, $3, $4, $5, $6, NOW(), NOW())`

	_, err := r.db.ExecContext(ctx, query,
		sub.ID,
		sub.ServiceName,
		sub.Price,
		sub.UserID,
		sub.StartDate,
		sub.EndDate,
	)
	return err
}

func (r *SubscriptionRepository) Update(ctx context.Context, sub *models.Subscription) error {
	query := `UPDATE subscriptions
			  SET service_name = $1, price = $2, user_id = $3, start_date = $4, end_date = $5, updated_at = NOW()
			  WHERE id = $6`

	result, err := r.db.ExecContext(ctx, query,
		sub.ServiceName,
		sub.Price,
		sub.UserID,
		sub.StartDate,
		sub.EndDate,
		sub.ID,
	)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *SubscriptionRepository) Get(ctx context.Context, id string) (*models.Subscription, error) {
	query := `SELECT id, service_name, price, user_id, start_date, end_date, created_at, updated_at
              FROM subscriptions
              WHERE id = $1`

	row := r.db.QueryRowContext(ctx, query, id)
	sub, err := scanSubscription(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return sub, nil
}

func (r *SubscriptionRepository) GetList(ctx context.Context, limit, offset int) ([]*models.Subscription, int, error) {
	var total int
	err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM subscriptions`).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	query := `SELECT id, service_name, price, user_id, start_date, end_date, created_at, updated_at
              FROM subscriptions
              ORDER BY created_at DESC
              LIMIT $1 OFFSET $2`

	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var list []*models.Subscription
	for rows.Next() {
		sub, err := scanSubscription(rows)
		if err != nil {
			return nil, 0, err
		}
		list = append(list, sub)
	}

	return list, total, rows.Err()
}

func (r *SubscriptionRepository) Delete(ctx context.Context, id string) (bool, error) {
	result, err := r.db.ExecContext(ctx, `DELETE FROM subscriptions WHERE id = $1`, id)
	if err != nil {
		return false, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return false, err
	}
	return rowsAffected > 0, nil
}

// GetByPeriod возвращает подписки, пересекающиеся с заданным периодом.
// Расчёт суммы с учётом месяцев выполняется в service-слое.
func (r *SubscriptionRepository) GetByPeriod(ctx context.Context, startDate, endDate, userID, serviceName string) ([]*models.Subscription, error) {
	query := `SELECT id, service_name, price, user_id, start_date, end_date, created_at, updated_at
              FROM subscriptions
              WHERE start_date <= $2
                AND (end_date IS NULL OR end_date >= $1)`

	args := []interface{}{startDate, endDate}

	if userID != "" {
		query += ` AND user_id = $3`
		args = append(args, userID)
	}
	if serviceName != "" {
		placeholder := `$3`
		if userID != "" {
			placeholder = `$4`
		}
		query += ` AND service_name = ` + placeholder
		args = append(args, serviceName)
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*models.Subscription
	for rows.Next() {
		sub, err := scanSubscription(rows)
		if err != nil {
			return nil, err
		}
		list = append(list, sub)
	}
	return list, rows.Err()
}

type rowScanner interface {
	Scan(dest ...any) error
}

func scanSubscription(s rowScanner) (*models.Subscription, error) {
	var sub models.Subscription
	err := s.Scan(
		&sub.ID,
		&sub.ServiceName,
		&sub.Price,
		&sub.UserID,
		&sub.StartDate,
		&sub.EndDate,
		&sub.CreatedAt,
		&sub.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &sub, nil
}
