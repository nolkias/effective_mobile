package router

import (
	"effective_mobile/internal/handlers"
	"github.com/gin-gonic/gin"
)

func SetupRouter(h *handlers.SubscriptionHandler) *gin.Engine {
	r := gin.Default()

	api := r.Group("/api/v1")

	{
		// CRUDL C ТЗ
		api.GET("/subscriptions/list", h.GetList)
		api.GET("/subscriptions/total-cost", h.TotalCost)
		api.GET("/subscriptions/:id", h.Get)

		api.POST("/subscriptions", h.Create)
		api.PUT("/subscriptions/:id", h.Update)
		api.DELETE("/subscriptions/:id", h.Delete)

	}

	return r
}
