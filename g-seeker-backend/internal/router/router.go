package router

import (
	"github.com/gin-gonic/gin"

	"g-seeker-backend/internal/client"
	"g-seeker-backend/internal/handler"
	"g-seeker-backend/internal/service"
)

func RegisterRoutes(r *gin.Engine) {
	githubClient := client.NewGitHubClient()
	githubSearchService := service.NewGitHubSearchService(githubClient)
	recommendService := service.NewRecommendService(githubSearchService)
	searchHandler := handler.NewSearchHandler(recommendService)

	api := r.Group("/api")
	{
		api.GET("/search", searchHandler.Search)
	}
}
