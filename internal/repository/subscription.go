package repository

import (
	"database/sql"
	"effective_mobile/internal/models"
	"errors"
)

// Scannable Для универсальности, чтобы не плодить для двух методов
type Scannable interface {
	Scan(dest ...any) error
}

// SubscriptionRepository В более серьзных проектах я бы делал через ORM, но тут набросл запросы
type SubscriptionRepository struct {
	db *sql.DB
}

func NewSubscriptionRepository(db *sql.DB) *SubscriptionRepository {
	return &SubscriptionRepository{db: db}
}

func (r *SubscriptionRepository) Create(sub *models.Subscription) error {
	query := `INSERT INTO subscriptions (id, service_name, price, user_id, start_date, end_date, created_at, updated_at)
              VALUES ($1, $2, $3, $4, $5, $6, NOW(), NOW())`

	_, err := r.db.Exec(
		query,
		sub.ID,
		sub.ServiceName,
		sub.Price,
		sub.UserID,
		sub.StartDate,
		sub.EndDate,
	)
	return err
}

func (r *SubscriptionRepository) Update(sub *models.Subscription) error {
	query := `UPDATE subscriptions 
			  SET service_name = $1, price = $2, user_id = $3, start_date = $4, end_date = $5, updated_at = NOW()
			  WHERE id = $6`

	result, err := r.db.Exec(
		query,
		sub.ServiceName,
		sub.Price,
		sub.UserID,
		sub.StartDate,
		sub.EndDate,
		sub.ID,
	)

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}
	return err
}

func (r *SubscriptionRepository) Get(id string) (*models.Subscription, error) {
	query := `SELECT id, service_name, price, user_id, start_date, end_date, created_at, updated_at 
              FROM subscriptions
              WHERE id = $1`
	row := r.db.QueryRow(query, id)

	result, err := toDTO(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return result, nil
}

func (r *SubscriptionRepository) GetList() ([]*models.Subscription, error) {
	list := make([]*models.Subscription, 0)
	query := `SELECT id, service_name, price, user_id, start_date, end_date, created_at, updated_at FROM subscriptions`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		sub, err := toDTO(rows)
		if err != nil {
			return nil, err
		}

		list = append(list, sub)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return list, nil
}

func (r *SubscriptionRepository) Delete(id string) error {
	query := `DELETE FROM subscriptions WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}

func (r *SubscriptionRepository) TotalCost(startDate, endDate, userID, serviceName string) (int, error) {
	query := `SELECT COALESCE(SUM(price), 0)
              FROM subscriptions
              WHERE (end_date IS NULL OR end_date >= $1)
                AND start_date <= $2
                AND user_id = $3
                AND service_name = $4`

	var total int
	err := r.db.QueryRow(query, startDate, endDate, userID, serviceName).Scan(&total)
	return total, err
}

func toDTO(scanner Scannable) (*models.Subscription, error) {
	var sub models.Subscription
	err := scanner.Scan(
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
