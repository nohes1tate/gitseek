package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"g-seeker-backend/internal/dto"
	"g-seeker-backend/internal/service"
	"g-seeker-backend/pkg/response"
)

type SearchHandler struct {
	recommendService service.RecommendService
}

func NewSearchHandler(recommendService service.RecommendService) *SearchHandler {
	return &SearchHandler{
		recommendService: recommendService,
	}
}

func (h *SearchHandler) Search(c *gin.Context) {
	var req dto.SearchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid request body")
		return
	}

	repos, err := h.recommendService.Search(c.Request.Context(), req.Query)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "failed to search repositories")
		return
	}

	items := make([]dto.SearchItem, 0, len(repos))
	for _, repo := range repos {
		items = append(items, dto.SearchItem{
			Name:        repo.Name,
			Owner:       repo.Owner,
			URL:         repo.URL,
			Stars:       repo.Stars,
			Description: repo.Description,
			Reason:      repo.Reason,
		})
	}

	response.Success(c, dto.SearchResponse{
		Items: items,
	})
}
