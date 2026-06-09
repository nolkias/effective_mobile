package router

import (
	_ "effective_mobile/docs"
	"effective_mobile/internal/handlers"
	"effective_mobile/internal/middleware"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRouter(h *handlers.SubscriptionHandler) *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middleware.Logger())

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := r.Group("/api/v1")
	{
		api.GET("/subscriptions/list", h.GetList)
		api.GET("/subscriptions/total-cost", h.TotalCost)
		api.GET("/subscriptions/:id", h.Get)
		api.POST("/subscriptions", h.Create)
		api.PUT("/subscriptions/:id", h.Update)
		api.DELETE("/subscriptions/:id", h.Delete)
	}

	return r
}
