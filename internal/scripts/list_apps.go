package scripts

import (
	"context"
	"monitoring/internal/domain"
	"strings"
)

type ListAppsReq struct {
	UserID     string `json:"userId"`
	Page       int    `json:"page"`
	Limit      int    `json:"limit"`
	SortOrder  string `json:"sortOrder"`
	SearchTerm string `json:"searchTerm"`
}

type ListAppsResp struct {
	Apps []domain.App `json:"apps"`
}

type ListAppsScript struct {
	appRepo domain.AppRepo
}

func NewListAppsScript(appRepo domain.AppRepo) *ListAppsScript {
	return &ListAppsScript{appRepo: appRepo}
}

func (s *ListAppsScript) Exec(ctx context.Context, req ListAppsReq) (*ListAppsResp, error) {
	userID, err := domain.NewID(req.UserID)
	if err != nil {
		return nil, err
	}

	filters := []domain.Filter{
		domain.NewFilter("userId", domain.Equals, userID),
	}

	if strings.TrimSpace(req.SearchTerm) != "" {
		filters = append(filters, domain.NewFilter("name", domain.Like, req.SearchTerm))
	}

	criteria := domain.NewCriteria(
		filters,
		domain.NewPagination(req.Limit, (req.Page-1)*req.Limit),
		domain.NewSort("createdAt", domain.SortOrder(req.SortOrder)),
	)

	apps, err := s.appRepo.ListApps(ctx, criteria)
	if err != nil {
		return nil, err
	}

	return &ListAppsResp{Apps: apps}, nil
}
