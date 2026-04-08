package service

import (
	"context"
	"fmt"
	"strings"

	"g-seeker-backend/internal/client"
	"g-seeker-backend/internal/model"
)

type GitHubSearchService struct {
	githubClient *client.GitHubClient
}

func NewGitHubSearchService(githubClient *client.GitHubClient) *GitHubSearchService {
	return &GitHubSearchService{
		githubClient: githubClient,
	}
}

func (s *GitHubSearchService) Search(ctx context.Context, userQuery string, limit int) ([]model.Repo, error) {
	if limit <= 0 {
		limit = 5
	}

	query := buildGitHubQuery(userQuery)

	resp, err := s.githubClient.SearchRepositories(ctx, query, 1, limit)
	if err != nil {
		return nil, err
	}

	repos := make([]model.Repo, 0, len(resp.Items))
	for _, item := range resp.Items {
		repos = append(repos, model.Repo{
			Name:        item.Name,
			Owner:       item.Owner.Login,
			URL:         item.HTMLURL,
			Stars:       item.StargazersCount,
			Description: item.Description,
			Reason:      buildReason(item),
		})
	}

	return repos, nil
}

func buildGitHubQuery(userQuery string) string {
	q := strings.TrimSpace(userQuery)
	if q == "" {
		return "stars:>100"
	}

	keywords := tokenizeSearchText(q)
	if len(keywords) == 0 {
		return "stars:>100"
	}

	parts := make([]string, 0, len(keywords)+2)
	parts = append(parts, keywords...)

	// 关键词少时适当提高 stars 门槛，避免泛召回
	switch {
	case len(keywords) <= 2:
		parts = append(parts, "stars:>80")
	case len(keywords) <= 4:
		parts = append(parts, "stars:>30")
	default:
		parts = append(parts, "stars:>10")
	}

	return strings.Join(parts, " ")
}

func buildReason(item client.GitHubRepoItem) string {
	if item.StargazersCount >= 10000 {
		return fmt.Sprintf("社区热度较高，当前 Star 数为 %d。", item.StargazersCount)
	}
	if item.Description != "" {
		return "项目描述较完整，与当前搜索需求有一定匹配度。"
	}
	return "该项目命中了搜索关键词，可作为候选仓库进一步分析。"
}
