package scripts

import (
	"context"
	"monitoring/internal/persistence"
	"strings"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ListAppsReq struct {
	UserID     string `json:"userId"`
	Page       int    `json:"page"`
	Limit      int    `json:"limit"`
	SortOrder  string `json:"sortOrder"`
	SearchTerm string `json:"searchTerm"`
}

type ListAppsResp struct {
	Apps []persistence.App `json:"apps"`
}

type ListAppsScript struct {
	db *mongo.Database
}

func NewListAppsScript(db *mongo.Database) *ListAppsScript {
	return &ListAppsScript{db: db}
}

func (s *ListAppsScript) Exec(ctx context.Context, req ListAppsReq) (*ListAppsResp, error) {
	userID, err := primitive.ObjectIDFromHex(req.UserID)
	if err != nil {
		return nil, err
	}

	filters := []persistence.Filter{
		persistence.NewFilter("userId", persistence.Equals, userID),
	}

	if strings.TrimSpace(req.SearchTerm) != "" {
		filters = append(filters, persistence.NewFilter("name", persistence.Like, req.SearchTerm))
	}

	criteria := persistence.NewCriteria(
		filters,
		persistence.NewPagination(req.Limit, (req.Page-1)*req.Limit),
		persistence.NewSort("createdAt", persistence.SortOrder(req.SortOrder)),
	)

	apps, err := persistence.ListApps(ctx, s.db, criteria)
	if err != nil {
		return nil, err
	}

	return &ListAppsResp{Apps: apps}, nil
}
