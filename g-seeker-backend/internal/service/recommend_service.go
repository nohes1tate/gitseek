package service

import (
	"context"
	"strings"

	"g-seeker-backend/internal/model"
)

type RecommendService interface {
	Search(ctx context.Context, query string) ([]model.Repo, error)
}

type recommendService struct{}

func NewRecommendService() RecommendService {
	return &recommendService{}
}

func (s *recommendService) Search(ctx context.Context, query string) ([]model.Repo, error) {
	lowerQuery := strings.ToLower(query)

	items := []model.Repo{
		{
			Name:        "casbin",
			Owner:       "casbin",
			URL:         "https://github.com/casbin/casbin",
			Stars:       18000,
			Description: "An authorization library supporting ACL, RBAC, and ABAC models.",
			Reason:      buildReason(lowerQuery, "A good fit for backend authorization scenarios."),
		},
		{
			Name:        "keto",
			Owner:       "ory",
			URL:         "https://github.com/ory/keto",
			Stars:       4200,
			Description: "A permission server inspired by Google Zanzibar for fine-grained access control.",
			Reason:      buildReason(lowerQuery, "Suitable when you need a more scalable permission model."),
		},
		{
			Name:        "permify",
			Owner:       "Permify",
			URL:         "https://github.com/Permify/permify",
			Stars:       5200,
			Description: "An open-source authorization service for building scalable access control.",
			Reason:      buildReason(lowerQuery, "A solid option if you prefer a standalone auth service."),
		},
	}

	return items, nil
}

func buildReason(query string, defaultReason string) string {
	switch {
	case strings.Contains(query, "go"):
		return "Matches your Go stack requirement. " + defaultReason
	case strings.Contains(query, "权限"), strings.Contains(query, "授权"):
		return "Matches your authorization requirement. " + defaultReason
	case strings.Contains(query, "中小团队"):
		return "Looks suitable for small to medium teams. " + defaultReason
	default:
		return defaultReason
	}
}
