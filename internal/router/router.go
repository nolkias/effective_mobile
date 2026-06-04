package router

import (
	_ "effective_mobile/docs"
	"effective_mobile/internal/handlers"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRouter(h *handlers.SubscriptionHandler) *gin.Engine {
	r := gin.Default()

	api := r.Group("/api/v1")

	{
		// Swagger UI
		r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

		// CRUDL C ТЗ
		// Даты возвращаются юез форматирования, так как обычно фронт решает как выводить, так что это не баг)
		api.GET("/subscriptions/list", h.GetList)
		api.GET("/subscriptions/total-cost", h.TotalCost)
		api.GET("/subscriptions/:id", h.Get)

		api.POST("/subscriptions", h.Create)
		api.PUT("/subscriptions/:id", h.Update)
		api.DELETE("/subscriptions/:id", h.Delete)

	}

	return r
}
