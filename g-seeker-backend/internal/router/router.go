package router

import (
	"github.com/gin-gonic/gin"

	"g-seeker-backend/internal/handler"
	"g-seeker-backend/internal/service"
)

func RegisterRoutes(r *gin.Engine) {
	api := r.Group("/api")

	recommendService := service.NewRecommendService()
	searchHandler := handler.NewSearchHandler(recommendService)

	api.POST("/search", searchHandler.Search)
}
