package service

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"g-seeker-backend/internal/llm"
	"g-seeker-backend/internal/model"
	"g-seeker-backend/internal/prompt"

	"github.com/cloudwego/eino/schema"
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
	llmClient           llm.Client
}

func NewRecommendService(
	githubSearchService *GitHubSearchService,
	rewriteService *QueryRewriteService,
	llmClient llm.Client,
) *RecommendService {
	return &RecommendService{
		githubSearchService: githubSearchService,
		rewriteService:      rewriteService,
		llmClient:           llmClient,
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
		ranked[i].Reason = s.buildRecommendationReasonWithFallback(
			ctx,
			query,
			rewriteResult.Query,
			ranked[i],
		)
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
	nameText := strings.ToLower(repo.Name)
	descText := strings.ToLower(repo.Description)
	fullText := nameText + " " + descText

	matchCount := 0
	nameMatchCount := 0

	for _, kw := range keywords {
		if strings.Contains(fullText, kw) {
			matchCount++
		}
		if strings.Contains(nameText, kw) {
			nameMatchCount++
		}
	}

	score := float64(matchCount) * 8
	score += float64(nameMatchCount) * 6

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
		score += 1.5
	}

	lowerName := strings.ToLower(repo.Name)
	lowerDesc := strings.ToLower(repo.Description)

	if strings.Contains(lowerName, "awesome") {
		score -= 3
	}
	if strings.Contains(lowerName, "example") || strings.Contains(lowerDesc, "example") {
		score -= 2
	}
	if strings.Contains(lowerName, "demo") || strings.Contains(lowerDesc, "demo") {
		score -= 2
	}
	if strings.Contains(lowerName, "tutorial") || strings.Contains(lowerDesc, "tutorial") {
		score -= 2
	}

	return score
}

func (s *RecommendService) buildRecommendationReasonWithFallback(
	ctx context.Context,
	originalQuery string,
	rewrittenQuery string,
	repo model.Repo,
) string {
	if s.llmClient != nil {
		reason, err := s.generateRecommendationReason(ctx, originalQuery, rewrittenQuery, repo)
		if err == nil && strings.TrimSpace(reason) != "" {
			return reason
		}
	}

	return buildRecommendationReasonFallback(originalQuery, rewrittenQuery, repo)
}

func (s *RecommendService) generateRecommendationReason(
	ctx context.Context,
	originalQuery string,
	rewrittenQuery string,
	repo model.Repo,
) (string, error) {
	messages := []*schema.Message{
		{
			Role:    schema.System,
			Content: prompt.RepoReasonSystemPrompt,
		},
		{
			Role: schema.User,
			Content: prompt.BuildRepoReasonUserPrompt(
				originalQuery,
				rewrittenQuery,
				repo.Name,
				repo.Owner,
				repo.Description,
				repo.Stars,
				repo.URL,
			),
		},
	}

	return s.llmClient.Generate(ctx, messages)
}

// func sanitizeReason(s string) string {
// 	s = strings.TrimSpace(s)
// 	s = strings.Trim(s, `"'`)
// 	s = strings.ReplaceAll(s, "\n", "")
// 	s = strings.Join(strings.Fields(s), " ")
// 	if len(s) > 120 {
// 		s = s[:120]
// 		s = strings.TrimSpace(s)
// 	}
// 	return s
// }

func buildRecommendationReasonFallback(originalQuery, rewrittenQuery string, repo model.Repo) string {
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
