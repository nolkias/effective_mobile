package handlers

import (
	"effective_mobile/internal/models"
	"effective_mobile/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"strconv"
)

type SubscriptionHandler struct {
	service service.SubscriptionSvc
}

func NewSubscriptionHandler(s service.SubscriptionSvc) *SubscriptionHandler {
	return &SubscriptionHandler{service: s}
}

// Create создаёт новую подписку
// @Summary Создать подписку
// @Description Создаёт новую подписку для пользователя
// @Tags Subscriptions
// @Accept json
// @Produce json
// @Param request body models.Subscription true "Данные подписки"
// @Success 201 {object} models.Subscription
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /subscriptions [post]
func (h *SubscriptionHandler) Create(c *gin.Context) {
	var req models.Subscription

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	req.ID = uuid.New()

	if err := h.service.Create(c.Request.Context(), &req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, req)
}

// Get возвращает подписку по ID
// @Summary Получить подписку
// @Description Возвращает подписку по её ID
// @Tags Subscriptions
// @Produce json
// @Param id path string true "ID подписки (UUID)"
// @Success 200 {object} models.Subscription
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /subscriptions/{id} [get]
func (h *SubscriptionHandler) Get(c *gin.Context) {
	id := c.Param("id")

	sub, err := h.service.Get(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if sub == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "subscription not found"})
		return
	}

	c.JSON(http.StatusOK, sub)
}

// GetList возвращает список всех подписок
// @Summary Список подписок
// @Description Возвращает список всех подписок
// @Tags Subscriptions
// @Produce json
// @Param limit query int false "Лимит (по умолчанию 20, максимум 100)"
// @Param offset query int false "Смещение (по умолчанию 0)"
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /subscriptions/list [get]
func (h *SubscriptionHandler) GetList(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "20")
	offsetStr := c.DefaultQuery("offset", "0")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		offset = 0
	}

	subs, total, err := h.service.GetList(c.Request.Context(), limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": subs,
		"pagination": gin.H{
			"limit":  limit,
			"offset": offset,
			"total":  total,
		},
	})
}

// Update обновляет существующую подписку
// @Summary Обновить подписку
// @Description Обновляет поля существующей подписки
// @Tags Subscriptions
// @Accept json
// @Produce json
// @Param id path string true "ID подписки (UUID)"
// @Param request body models.Subscription true "Данные для обновления"
// @Success 200 {object} models.Subscription
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /subscriptions/{id} [put]
func (h *SubscriptionHandler) Update(c *gin.Context) {
	id := c.Param("id")
	ctx := c.Request.Context()

	var req models.Subscription
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	existing, err := h.service.Get(ctx, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if existing == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "subscription not found"})
		return
	}

	if req.ServiceName != "" {
		existing.ServiceName = req.ServiceName
	}
	if req.Price != 0 {
		existing.Price = req.Price
	}
	if req.UserID != uuid.Nil {
		existing.UserID = req.UserID
	}
	if req.StartDate != "" {
		existing.StartDate = req.StartDate
	}
	if req.EndDate != nil {
		existing.EndDate = req.EndDate
	}

	if err := h.service.Update(ctx, existing); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, existing)
}

// Delete удаляет подписку
// @Summary Удалить подписку
// @Description Удаляет подписку по ID
// @Tags Subscriptions
// @Produce json
// @Param id path string true "ID подписки (UUID)"
// @Success 204 "No Content"
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /subscriptions/{id} [delete]
func (h *SubscriptionHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	found, err := h.service.Delete(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if !found {
		c.JSON(http.StatusNotFound, gin.H{"error": "subscription not found"})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

// TotalCost подсчитывает стоимость подписок за период
// @Summary Подсчёт стоимости
// @Description Подсчитывает суммарную стоимость подписок за указанный период с учётом количества месяцев
// @Tags Subscriptions
// @Produce json
// @Param start_date query string true "Начало периода (MM-YYYY)"
// @Param end_date query string true "Конец периода (MM-YYYY)"
// @Param user_id query string false "ID пользователя (UUID)"
// @Param service_name query string false "Название сервиса"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /subscriptions/total-cost [get]
func (h *SubscriptionHandler) TotalCost(c *gin.Context) {
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")
	userID := c.Query("user_id")
	serviceName := c.Query("service_name")

	if startDate == "" || endDate == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "start_date and end_date are required"})
		return
	}

	total, err := h.service.TotalCost(c.Request.Context(), startDate, endDate, userID, serviceName)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"total_cost": total,
		"period": gin.H{
			"start_date": startDate,
			"end_date":   endDate,
		},
	})
}
