package service

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"g-seeker-backend/internal/model"
)

type RecommendResult struct {
	OriginalQuery  string       `json:"original_query"`
	RewrittenQuery string       `json:"rewritten_query"`
	CandidateCount int          `json:"candidate_count"`
	Items          []model.Repo `json:"items"`
}

type RecommendService struct {
	githubSearchService *GitHubSearchService
	rewriteService      *QueryRewriteService
}

func NewRecommendService(
	githubSearchService *GitHubSearchService,
	rewriteService *QueryRewriteService,
) *RecommendService {
	return &RecommendService{
		githubSearchService: githubSearchService,
		rewriteService:      rewriteService,
	}
}

func (s *RecommendService) Recommend(ctx context.Context, query string, limit int) (*RecommendResult, error) {
	if limit <= 0 {
		limit = 5
	}
	if limit > 10 {
		limit = 10
	}

	rewriteResult, err := s.rewriteService.Rewrite(ctx, query)
	if err != nil {
		return nil, err
	}

	candidateLimit := limit * 3
	if candidateLimit < 10 {
		candidateLimit = 10
	}
	if candidateLimit > 30 {
		candidateLimit = 30
	}

	repos, err := s.githubSearchService.Search(ctx, rewriteResult.Query, candidateLimit)
	if err != nil {
		return nil, err
	}

	ranked := rankRepos(query, rewriteResult.Query, repos)

	if len(ranked) > limit {
		ranked = ranked[:limit]
	}

	for i := range ranked {
		ranked[i].Reason = buildRecommendationReason(query, rewriteResult.Query, ranked[i])
	}

	return &RecommendResult{
		OriginalQuery:  query,
		RewrittenQuery: rewriteResult.Query,
		CandidateCount: len(repos),
		Items:          ranked,
	}, nil
}

type scoredRepo struct {
	repo  model.Repo
	score float64
}

func rankRepos(originalQuery, rewrittenQuery string, repos []model.Repo) []model.Repo {
	keywords := uniqueStrings(append(
		tokenizeSearchText(originalQuery),
		tokenizeSearchText(rewrittenQuery)...,
	))

	scored := make([]scoredRepo, 0, len(repos))
	for _, repo := range repos {
		scored = append(scored, scoredRepo{
			repo:  repo,
			score: calculateRepoScore(repo, keywords),
		})
	}

	sort.SliceStable(scored, func(i, j int) bool {
		if scored[i].score == scored[j].score {
			if scored[i].repo.Stars == scored[j].repo.Stars {
				return scored[i].repo.Name < scored[j].repo.Name
			}
			return scored[i].repo.Stars > scored[j].repo.Stars
		}
		return scored[i].score > scored[j].score
	})

	out := make([]model.Repo, 0, len(scored))
	for _, item := range scored {
		out = append(out, item.repo)
	}
	return out
}

func calculateRepoScore(repo model.Repo, keywords []string) float64 {
	text := strings.ToLower(repo.Name + " " + repo.Description)

	matchCount := 0
	for _, kw := range keywords {
		if strings.Contains(text, kw) {
			matchCount++
		}
	}

	score := float64(matchCount) * 10

	switch {
	case repo.Stars >= 50000:
		score += 12
	case repo.Stars >= 10000:
		score += 9
	case repo.Stars >= 3000:
		score += 6
	case repo.Stars >= 1000:
		score += 4
	case repo.Stars >= 200:
		score += 2
	}

	if repo.Description != "" {
		score += 1
	}

	if strings.Contains(strings.ToLower(repo.Name), "awesome") {
		score -= 2
	}

	return score
}

func buildRecommendationReason(originalQuery, rewrittenQuery string, repo model.Repo) string {
	keywords := uniqueStrings(append(
		tokenizeSearchText(originalQuery),
		tokenizeSearchText(rewrittenQuery)...,
	))

	matched := matchedKeywords(repo, keywords)
	switch {
	case len(matched) >= 2 && repo.Stars >= 1000:
		return fmt.Sprintf(
			"与查询关键词「%s」匹配度较高，且社区热度较好（%d Stars），适合作为优先候选。",
			strings.Join(matched[:minInt(2, len(matched))], "、"),
			repo.Stars,
		)
	case len(matched) >= 1:
		return fmt.Sprintf(
			"命中了关键词「%s」，仓库描述与当前需求较接近，可作为重点评估对象。",
			strings.Join(matched[:minInt(2, len(matched))], "、"),
		)
	case repo.Stars >= 10000:
		return fmt.Sprintf(
			"社区热度较高（%d Stars），在同类候选中具备较强参考价值。",
			repo.Stars,
		)
	case repo.Description != "":
		return "仓库描述较完整，和当前搜索意图存在一定相关性，建议进一步查看 README 与示例代码。"
	default:
		return "该仓库出现在当前搜索结果中，可作为候选项目进一步分析。"
	}
}

func matchedKeywords(repo model.Repo, keywords []string) []string {
	text := strings.ToLower(repo.Name + " " + repo.Description)
	out := make([]string, 0, 3)
	for _, kw := range keywords {
		if strings.Contains(text, kw) {
			out = append(out, kw)
		}
		if len(out) >= 3 {
			break
		}
	}
	return out
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}
