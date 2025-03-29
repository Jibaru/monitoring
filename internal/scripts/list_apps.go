package scripts

import (
	"context"
	"monitoring/internal/persistence"

	"go.mongodb.org/mongo-driver/mongo"
)

type ListAppsReq struct {
	UserID    string `json:"userId"`
	Page      int    `json:"page"`
	Limit     int    `json:"limit"`
	SortOrder string `json:"sortOrder"`
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
	criteria := persistence.NewCriteria(
		[]persistence.Filter{
			persistence.NewFilter("userId", persistence.Equals, req.UserID),
		},
		persistence.NewPagination(req.Limit, (req.Page-1)*req.Limit),
		persistence.NewSort("createdAt", persistence.SortOrder(req.SortOrder)),
	)

	apps, err := persistence.ListAppsPaginated(ctx, s.db, criteria)
	if err != nil {
		return nil, err
	}

	return &ListAppsResp{Apps: apps}, nil
}
