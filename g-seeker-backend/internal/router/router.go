package router

import (
	"log"

	"github.com/gin-gonic/gin"

	"g-seeker-backend/internal/client"
	"g-seeker-backend/internal/handler"
	"g-seeker-backend/internal/llm"
	"g-seeker-backend/internal/service"
)

func RegisterRoutes(r *gin.Engine) {
	githubClient := client.NewGitHubClient()
	githubSearchService := service.NewGitHubSearchService(githubClient)

	var llmClient llm.Client
	lc, err := llm.NewLLMClient()
	if err != nil {
		log.Printf("warning: llm client init failed, fallback to rule-based rewrite only: %v", err)
	} else {
		llmClient = lc
		log.Printf("info: llm client initialized successfully")
	}

	rewriteService := service.NewQueryRewriteService(llmClient)
	recommendService := service.NewRecommendService(githubSearchService, rewriteService)
	searchHandler := handler.NewSearchHandler(recommendService)

	api := r.Group("/api")
	{
		api.GET("/search", searchHandler.Search)
	}
}
