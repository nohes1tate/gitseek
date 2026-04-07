package client

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

type GitHubClient struct {
	baseURL    string
	token      string
	httpClient *http.Client
}

func NewGitHubClient() *GitHubClient {
	baseURL := os.Getenv("GITHUB_API_BASE")
	if baseURL == "" {
		baseURL = "https://api.github.com"
	}

	return &GitHubClient{
		baseURL: strings.TrimRight(baseURL, "/"),
		token:   os.Getenv("GITHUB_TOKEN"),
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

type SearchRepositoriesResponse struct {
	TotalCount int              `json:"total_count"`
	Items      []GitHubRepoItem `json:"items"`
}

type GitHubRepoItem struct {
	Name            string      `json:"name"`
	FullName        string      `json:"full_name"`
	HTMLURL         string      `json:"html_url"`
	Description     string      `json:"description"`
	StargazersCount int         `json:"stargazers_count"`
	Owner           GitHubOwner `json:"owner"`
}

type GitHubOwner struct {
	Login string `json:"login"`
}

func (c *GitHubClient) SearchRepositories(ctx context.Context, query string, page, perPage int) (*SearchRepositoriesResponse, error) {
	if strings.TrimSpace(query) == "" {
		return nil, fmt.Errorf("query is empty")
	}
	if page <= 0 {
		page = 1
	}
	if perPage <= 0 || perPage > 100 {
		perPage = 10
	}

	u, err := url.Parse(c.baseURL + "/search/repositories")
	if err != nil {
		return nil, err
	}

	params := u.Query()
	params.Set("q", query)
	params.Set("page", fmt.Sprintf("%d", page))
	params.Set("per_page", fmt.Sprintf("%d", perPage))
	params.Set("sort", "stars")
	params.Set("order", "desc")
	u.RawQuery = params.Encode()

	fmt.Printf("github baseURL=%s\n", c.baseURL)
	fmt.Printf("github token loaded=%v, len=%d\n", c.token != "", len(c.token))

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")
	req.Header.Set("User-Agent", "g-seeker-backend")

	if c.token != "" {
		req.Header.Set("Authorization", "Bearer "+c.token)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 300 {
		return nil, fmt.Errorf("github search failed, status=%d, body=%s", resp.StatusCode, string(body))
	}

	var result SearchRepositoriesResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	return &result, nil
}
