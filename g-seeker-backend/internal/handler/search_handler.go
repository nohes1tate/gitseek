package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"g-seeker-backend/internal/service"
)

type SearchHandler struct {
	recommendService *service.RecommendService
}

func NewSearchHandler(recommendService *service.RecommendService) *SearchHandler {
	return &SearchHandler{
		recommendService: recommendService,
	}
}

func (h *SearchHandler) Search(c *gin.Context) {
	query := c.Query("query")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "query is required",
		})
		return
	}

	limit := 5
	if limitStr := c.Query("limit"); limitStr != "" {
		if v, err := strconv.Atoi(limitStr); err == nil && v > 0 {
			limit = v
		}
	}

	repos, err := h.recommendService.Recommend(c.Request.Context(), query, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "search failed",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "success",
		"data":    repos,
	})
}
