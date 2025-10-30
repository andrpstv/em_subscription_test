package service

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"em_subscription_test/internal/repository"
	"em_subscription_test/models"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type SubscriptionService interface {
	Create(req *models.SubscriptionCreate) (*models.Subscription, error)
	GetByID(id uuid.UUID) (*models.Subscription, error)
	List(userID *uuid.UUID, serviceName *string) ([]models.Subscription, error)
	Update(id uuid.UUID, req *models.SubscriptionUpdate) (*models.Subscription, error)
	Delete(id uuid.UUID) error
	GetTotalCost(req *models.TotalCostRequest) (*models.TotalCostResponse, error)
}

type subscriptionService struct {
	repo   repository.SubscriptionRepository
	logger *logrus.Logger
}

func NewSubscriptionService(repo repository.SubscriptionRepository, logger *logrus.Logger) SubscriptionService {
	return &subscriptionService{repo: repo, logger: logger}
}

func (s *subscriptionService) Create(req *models.SubscriptionCreate) (*models.Subscription, error) {
	if !isValidDateFormat(req.StartDate) {
		return nil, fmt.Errorf("start_date must be in MM-YYYY format")
	}
	if req.EndDate != nil && !isValidDateFormat(*req.EndDate) {
		return nil, fmt.Errorf("end_date must be in MM-YYYY format")
	}

	subscription := &models.Subscription{
		ID:          uuid.New(),
		ServiceName: req.ServiceName,
		Price:       req.Price,
		UserID:      req.UserID,
		StartDate:   req.StartDate,
		EndDate:     req.EndDate,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	err := s.repo.Create(subscription)
	if err != nil {
		s.logger.WithError(err).Error("Failed to create subscription")
		return nil, err
	}

	s.logger.WithFields(logrus.Fields{
		"id":           subscription.ID,
		"service_name": subscription.ServiceName,
		"user_id":      subscription.UserID,
	}).Info("Subscription created")

	return subscription, nil
}

func (s *subscriptionService) GetByID(id uuid.UUID) (*models.Subscription, error) {
	subscription, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("not found")
		}
		s.logger.WithError(err).Error("Failed to get subscription")
		return nil, err
	}
	return subscription, nil
}

func (s *subscriptionService) List(userID *uuid.UUID, serviceName *string) ([]models.Subscription, error) {
	filters := make(map[string]interface{})
	if userID != nil {
		filters["user_id"] = *userID
	}
	if serviceName != nil {
		filters["service_name"] = *serviceName
	}

	subscriptions, err := s.repo.List(filters)
	if err != nil {
		s.logger.WithError(err).Error("Failed to list subscriptions")
		return nil, err
	}
	if subscriptions == nil {
		return []models.Subscription{}, nil
	}
	return subscriptions, nil
}

func (s *subscriptionService) Update(id uuid.UUID, req *models.SubscriptionUpdate) (*models.Subscription, error) {
	if req.StartDate != nil && !isValidDateFormat(*req.StartDate) {
		return nil, fmt.Errorf("start_date must be in MM-YYYY format")
	}
	if req.EndDate != nil && *req.EndDate != "" && !isValidDateFormat(*req.EndDate) {
		return nil, fmt.Errorf("end_date must be in MM-YYYY format")
	}

	existing, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if req.ServiceName != nil {
		existing.ServiceName = *req.ServiceName
	}
	if req.Price != nil {
		existing.Price = *req.Price
	}
	if req.UserID != nil {
		existing.UserID = *req.UserID
	}
	if req.StartDate != nil {
		existing.StartDate = *req.StartDate
	}
	if req.EndDate != nil {
		existing.EndDate = req.EndDate
	}
	existing.UpdatedAt = time.Now()

	err = s.repo.Update(existing)
	if err != nil {
		s.logger.WithError(err).Error("Failed to update subscription")
		return nil, err
	}

	s.logger.WithField("id", id).Info("Subscription updated")
	return existing, nil
}

func (s *subscriptionService) Delete(id uuid.UUID) error {
	err := s.repo.Delete(id)
	if err != nil {
		s.logger.WithError(err).Error("Failed to delete subscription")
		return err
	}
	s.logger.WithField("id", id).Info("Subscription deleted")
	return nil
}

func (s *subscriptionService) GetTotalCost(req *models.TotalCostRequest) (*models.TotalCostResponse, error) {
	if !isValidDateFormat(req.StartPeriod) || !isValidDateFormat(req.EndPeriod) {
		return nil, fmt.Errorf("start_period and end_period must be in MM-YYYY format")
	}

	startPeriod, err := parsePeriod(req.StartPeriod)
	if err != nil {
		return nil, fmt.Errorf("invalid start_period")
	}
	endPeriod, err := parsePeriod(req.EndPeriod)
	if err != nil {
		return nil, fmt.Errorf("invalid end_period")
	}

	if startPeriod.After(endPeriod) {
		return nil, fmt.Errorf("start_period must be before or equal to end_period")
	}

	filters := make(map[string]interface{})
	if req.UserID != nil {
		filters["user_id"] = *req.UserID
	}
	if req.ServiceName != nil {
		filters["service_name"] = *req.ServiceName
	}

	subscriptions, err := s.repo.List(filters)
	if err != nil {
		s.logger.WithError(err).Error("Failed to get subscriptions for total cost")
		return nil, err
	}

	totalCost := 0
	for _, sub := range subscriptions {
		subStart, _ := parsePeriod(sub.StartDate)
		if sub.EndDate == nil {
			if !subStart.Before(startPeriod) && !subStart.After(endPeriod) {
				totalCost += sub.Price
			}
		} else {
			subEnd, _ := parsePeriod(*sub.EndDate)
			overlapMonths := calculateOverlapMonths(startPeriod, endPeriod, subStart, subEnd)
			totalCost += sub.Price * overlapMonths
		}
	}

	response := &models.TotalCostResponse{TotalCost: totalCost}

	s.logger.WithFields(logrus.Fields{
		"start_period": req.StartPeriod,
		"end_period":   req.EndPeriod,
		"user_id":      req.UserID,
		"service_name": req.ServiceName,
		"total_cost":   totalCost,
	}).Info("Total cost calculated")

	return response, nil
}

func isValidDateFormat(date string) bool {
	parts := strings.Split(date, "-")
	if len(parts) != 2 {
		return false
	}
	month, err1 := strconv.Atoi(parts[0])
	year, err2 := strconv.Atoi(parts[1])
	if err1 != nil || err2 != nil {
		return false
	}
	return month >= 1 && month <= 12 && year >= 1900 && year <= 9999
}

func parsePeriod(period string) (time.Time, error) {
	parts := strings.Split(period, "-")
	month, _ := strconv.Atoi(parts[0])
	year, _ := strconv.Atoi(parts[1])
	return time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC), nil
}

func calculateOverlapMonths(periodStart, periodEnd, subStart, subEnd time.Time) int {
	overlapStart := maxTime(periodStart, subStart)
	overlapEnd := minTime(periodEnd, subEnd)

	if overlapStart.After(overlapEnd) {
		return 0
	}

	yearDiff := overlapEnd.Year() - overlapStart.Year()
	monthDiff := int(overlapEnd.Month()) - int(overlapStart.Month())
	return yearDiff*12 + monthDiff + 1
}

func maxTime(a, b time.Time) time.Time {
	if a.After(b) {
		return a
	}
	return b
}

func minTime(a, b time.Time) time.Time {
	if a.Before(b) {
		return a
	}
	return b
}
