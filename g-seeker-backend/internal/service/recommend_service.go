package service

import (
	"context"

	"g-seeker-backend/internal/model"
)

type RecommendService struct {
	githubSearchService *GitHubSearchService
}

func NewRecommendService(githubSearchService *GitHubSearchService) *RecommendService {
	return &RecommendService{
		githubSearchService: githubSearchService,
	}
}

func (s *RecommendService) Recommend(ctx context.Context, query string, limit int) ([]model.Repo, error) {
	return s.githubSearchService.Search(ctx, query, limit)
}
