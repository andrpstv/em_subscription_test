package repository

import (
	"fmt"

	"em_subscription_test/models"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type SubscriptionRepository interface {
	Create(subscription *models.Subscription) error
	GetByID(id uuid.UUID) (*models.Subscription, error)
	List(filters map[string]interface{}) ([]models.Subscription, error)
	Update(subscription *models.Subscription) error
	Delete(id uuid.UUID) error
}

type subscriptionRepository struct {
	db *sqlx.DB
}

func NewSubscriptionRepository(db *sqlx.DB) SubscriptionRepository {
	return &subscriptionRepository{db: db}
}

func (r *subscriptionRepository) Create(subscription *models.Subscription) error {
	query := `INSERT INTO subscriptions (id, service_name, price, user_id, start_date, end_date, created_at, updated_at)
	          VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	_, err := r.db.Exec(query, subscription.ID, subscription.ServiceName, subscription.Price,
		subscription.UserID, subscription.StartDate, subscription.EndDate,
		subscription.CreatedAt, subscription.UpdatedAt)
	return err
}

func (r *subscriptionRepository) GetByID(id uuid.UUID) (*models.Subscription, error) {
	var subscription models.Subscription
	query := `SELECT id, service_name, price, user_id, start_date, end_date, created_at, updated_at
	          FROM subscriptions WHERE id = $1`
	err := r.db.Get(&subscription, query, id)
	if err != nil {
		return nil, err
	}
	return &subscription, nil
}

func (r *subscriptionRepository) List(filters map[string]interface{}) ([]models.Subscription, error) {
	query := `SELECT id, service_name, price, user_id, start_date, end_date, created_at, updated_at FROM subscriptions WHERE 1=1`
	args := []interface{}{}
	argCount := 0

	if userID, ok := filters["user_id"]; ok && userID != nil {
		argCount++
		query += fmt.Sprintf(" AND user_id = $%d", argCount)
		args = append(args, userID)
	}

	if serviceName, ok := filters["service_name"]; ok && serviceName != nil {
		argCount++
		query += fmt.Sprintf(" AND service_name = $%d", argCount)
		args = append(args, serviceName)
	}

	var subscriptions []models.Subscription
	err := r.db.Select(&subscriptions, query, args...)
	return subscriptions, err
}

func (r *subscriptionRepository) Update(subscription *models.Subscription) error {
	query := `UPDATE subscriptions SET service_name = $1, price = $2, user_id = $3,
	          start_date = $4, end_date = $5, updated_at = $6 WHERE id = $7`
	_, err := r.db.Exec(query, subscription.ServiceName, subscription.Price, subscription.UserID,
		subscription.StartDate, subscription.EndDate, subscription.UpdatedAt, subscription.ID)
	return err
}

func (r *subscriptionRepository) Delete(id uuid.UUID) error {
	query := `DELETE FROM subscriptions WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}
