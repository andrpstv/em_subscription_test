package handlers

import (
	"database/sql"
	"errors"
	"net/http"

	"em_subscription_test/internal/service"
	"em_subscription_test/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type Handler struct {
	Service service.SubscriptionService
	Logger  *logrus.Logger
}

func NewHandler(svc service.SubscriptionService, logger *logrus.Logger) *Handler {
	return &Handler{
		Service: svc,
		Logger:  logger,
	}
}

// CreateSubscription creates a new subscription
// @Summary Create a new subscription
// @Description Create a new subscription with the provided details
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param subscription body models.SubscriptionCreate true "Subscription data"
// @Success 201 {object} models.Subscription
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /subscriptions [post]
func (h *Handler) CreateSubscription(c *gin.Context) {
	var req models.SubscriptionCreate
	if err := c.ShouldBindJSON(&req); err != nil {
		h.Logger.WithError(err).Error("Invalid request body")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	subscription, err := h.Service.Create(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, subscription)
}

// GetSubscription gets a subscription by ID
// @Summary Get a subscription by ID
// @Description Get a subscription by its ID
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param id path string true "Subscription ID"
// @Success 200 {object} models.Subscription
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /subscriptions/{id} [get]
func (h *Handler) GetSubscription(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		h.Logger.WithError(err).Error("Invalid subscription ID")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid subscription ID"})
		return
	}

	subscription, err := h.Service.GetByID(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Subscription not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get subscription"})
		return
	}

	c.JSON(http.StatusOK, subscription)
}

// ListSubscriptions lists all subscriptions with optional filters
// @Summary List subscriptions
// @Description List all subscriptions with optional filtering by user_id and service_name
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param user_id query string false "User ID"
// @Param service_name query string false "Service Name"
// @Success 200 {array} models.Subscription
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /subscriptions [get]
func (h *Handler) ListSubscriptions(c *gin.Context) {
	userIDStr := c.Query("user_id")
	serviceName := c.Query("service_name")

	var userID *uuid.UUID
	if userIDStr != "" {
		parsed, err := uuid.Parse(userIDStr)
		if err != nil {
			h.Logger.WithError(err).Error("Invalid user_id")
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user_id"})
			return
		}
		userID = &parsed
	}

	var svcName *string
	if serviceName != "" {
		svcName = &serviceName
	}

	subscriptions, err := h.Service.List(userID, svcName)
	if err != nil {
		h.Logger.WithError(err).Error("Failed to list subscriptions")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list subscriptions"})
		return
	}

	c.JSON(http.StatusOK, subscriptions)
}

// UpdateSubscription updates a subscription by ID
// @Summary Update a subscription
// @Description Update a subscription by its ID
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param id path string true "Subscription ID"
// @Param subscription body models.SubscriptionUpdate true "Updated subscription data"
// @Success 200 {object} models.Subscription
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /subscriptions/{id} [put]
func (h *Handler) UpdateSubscription(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		h.Logger.WithError(err).Error("Invalid subscription ID")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid subscription ID"})
		return
	}

	var req models.SubscriptionUpdate
	if err := c.ShouldBindJSON(&req); err != nil {
		h.Logger.WithError(err).Error("Invalid request body")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	subscription, err := h.Service.Update(id, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, subscription)
}

// DeleteSubscription deletes a subscription by ID
// @Summary Delete a subscription
// @Description Delete a subscription by its ID
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param id path string true "Subscription ID"
// @Success 204
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /subscriptions/{id} [delete]
func (h *Handler) DeleteSubscription(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		h.Logger.WithError(err).Error("Invalid subscription ID")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid subscription ID"})
		return
	}

	err = h.Service.Delete(id)
	if err != nil {
		h.Logger.WithError(err).Error("Failed to delete subscription")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete subscription"})
		return
	}

	c.Status(http.StatusNoContent)
}

// GetTotalCost calculates the total cost of subscriptions for a given period
// @Summary Get total cost of subscriptions
// @Description Calculate the total cost of subscriptions for a given period with optional filters
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param request body models.TotalCostRequest true "Total cost request"
// @Success 200 {object} models.TotalCostResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /subscriptions/total-cost [post]
func (h *Handler) GetTotalCost(c *gin.Context) {
	var req models.TotalCostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.Logger.WithError(err).Error("Invalid request body")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.Service.GetTotalCost(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}
