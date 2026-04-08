package service

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"unicode"

	"g-seeker-backend/internal/llm"

	"github.com/cloudwego/eino/schema"
)

type RewriteResult struct {
	OriginalQuery string
	Query         string
	Keywords      []string
}

type QueryRewriteService struct {
	llmClient llm.Client
}

func NewQueryRewriteService(llmClient llm.Client) *QueryRewriteService {
	return &QueryRewriteService{
		llmClient: llmClient,
	}
}

func (s *QueryRewriteService) Rewrite(ctx context.Context, userQuery string) (*RewriteResult, error) {
	userQuery = strings.TrimSpace(userQuery)
	if userQuery == "" {
		return nil, fmt.Errorf("query is empty")
	}

	fallback := ruleBasedRewrite(userQuery)
	result := &RewriteResult{
		OriginalQuery: userQuery,
		Query:         fallback,
		Keywords:      tokenizeSearchText(fallback),
	}

	fmt.Printf("[rewrite] original query: %s\n", userQuery)
	fmt.Printf("[rewrite] fallback query: %s\n", fallback)

	if s.llmClient == nil {
		return result, nil
	}

	systemPrompt := `You are a GitHub repository search query rewriter.
Convert a user's natural language software requirement into a concise GitHub search query.

Rules:
1. Output English only.
2. Output one single line only.
3. Keep it short and searchable.
4. Include language/framework/domain keywords when helpful.
5. Do not explain anything.
6. Do not wrap with quotes or markdown.
7. Prefer keyword style, not sentence style.`

	userPrompt := fmt.Sprintf("User requirement: %s", userQuery)

	fmt.Printf("[rewrite] calling llm...\n")

	messages := []*schema.Message{
		{
			Role:    schema.System,
			Content: systemPrompt,
		},
		{
			Role:    schema.User,
			Content: userPrompt,
		},
	}

	rewritten, err := s.llmClient.Generate(ctx, messages)

	// rewritten, err := s.llmClient.Generate(ctx, systemPrompt, userPrompt)
	if err != nil {
		fmt.Printf("[rewrite] llm error: %v\n", err)
		return result, nil
	}

	fmt.Printf("[rewrite] raw llm output: %s\n", rewritten)

	rewritten = sanitizeGeneratedQuery(rewritten)
	if rewritten == "" {
		return result, nil
	}

	result.Query = rewritten
	result.Keywords = tokenizeSearchText(rewritten)

	fmt.Printf("[rewrite] final rewritten query: %s\n", result.Query)
	return result, nil
}

func sanitizeGeneratedQuery(s string) string {
	s = strings.TrimSpace(s)
	s = strings.Trim(s, `"'`)
	s = strings.ReplaceAll(s, "\n", " ")
	s = strings.Join(strings.Fields(s), " ")

	if len(s) > 120 {
		s = s[:120]
		s = strings.TrimSpace(s)
	}
	return s
}

func ruleBasedRewrite(q string) string {
	raw := strings.ToLower(strings.TrimSpace(q))

	parts := make([]string, 0, 6)

	if hasAny(raw, "go", "golang") {
		parts = append(parts, "golang")
	} else if hasAny(raw, "python", "py") {
		parts = append(parts, "python")
	} else if hasAny(raw, "java", "spring") {
		parts = append(parts, "java")
	} else if hasAny(raw, "typescript", "ts", "node", "next.js", "nextjs") {
		parts = append(parts, "typescript")
	} else if hasAny(raw, "javascript", "js") {
		parts = append(parts, "javascript")
	} else if hasAny(raw, "rust") {
		parts = append(parts, "rust")
	}

	switch {
	case hasAny(raw, "权限", "鉴权", "认证", "授权", "rbac", "auth", "oauth", "access control", "authorization", "authentication"):
		parts = append(parts, "authorization", "authentication", "rbac")
	case hasAny(raw, "工作流", "workflow", "orchestration", "dag"):
		parts = append(parts, "workflow", "engine", "orchestration")
	case hasAny(raw, "消息队列", "mq", "message queue", "pubsub", "pub sub", "kafka", "rabbitmq"):
		parts = append(parts, "message", "queue", "client")
	case hasAny(raw, "日志", "logger", "logging", "log"):
		parts = append(parts, "logging", "logger")
	case hasAny(raw, "配置", "config", "configuration"):
		parts = append(parts, "config", "configuration")
	case hasAny(raw, "搜索", "search", "全文检索", "elasticsearch"):
		parts = append(parts, "search", "engine")
	case hasAny(raw, "爬虫", "crawler", "scraper", "scraping"):
		parts = append(parts, "crawler", "scraper")
	case hasAny(raw, "orm", "数据库", "database", "sql"):
		parts = append(parts, "orm", "database")
	case hasAny(raw, "缓存", "cache", "redis"):
		parts = append(parts, "cache", "redis")
	default:
		parts = append(parts, tokenizeSearchText(q)...)
	}

	parts = uniqueStrings(parts)
	if len(parts) == 0 {
		return q
	}

	return strings.Join(parts, " ")
}

func hasAny(s string, words ...string) bool {
	for _, w := range words {
		if strings.Contains(s, strings.ToLower(w)) {
			return true
		}
	}
	return false
}

var nonWordRegexp = regexp.MustCompile(`[^a-zA-Z0-9#+._-]+`)

func tokenizeSearchText(s string) []string {
	s = strings.ToLower(s)

	builder := strings.Builder{}
	for _, r := range s {
		if r > unicode.MaxASCII {
			builder.WriteRune(' ')
			continue
		}
		builder.WriteRune(r)
	}
	s = builder.String()

	s = nonWordRegexp.ReplaceAllString(s, " ")
	items := strings.Fields(s)

	stopwords := map[string]struct{}{
		"the": {}, "a": {}, "an": {}, "for": {}, "of": {}, "to": {},
		"and": {}, "or": {}, "with": {}, "in": {}, "on": {}, "by": {},
		"is": {}, "are": {}, "repo": {}, "github": {}, "library": {},
		"project": {}, "tool": {}, "tools": {}, "best": {}, "help": {},
	}

	out := make([]string, 0, len(items))
	for _, item := range items {
		if len(item) <= 1 {
			continue
		}
		if _, ok := stopwords[item]; ok {
			continue
		}
		out = append(out, item)
	}
	return uniqueStrings(out)
}

func uniqueStrings(items []string) []string {
	seen := make(map[string]struct{}, len(items))
	out := make([]string, 0, len(items))
	for _, item := range items {
		item = strings.TrimSpace(strings.ToLower(item))
		if item == "" {
			continue
		}
		if _, ok := seen[item]; ok {
			continue
		}
		seen[item] = struct{}{}
		out = append(out, item)
	}
	return out
}
