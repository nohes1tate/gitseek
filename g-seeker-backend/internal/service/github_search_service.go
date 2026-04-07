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

	// MVP 阶段先做简单 query rewrite
	return fmt.Sprintf("%s stars:>50", q)
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
