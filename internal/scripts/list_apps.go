package scripts

import (
	"context"
	"monitoring/internal/persistence"

	"go.mongodb.org/mongo-driver/mongo"
)

type ListAppsReq struct {
	UserID string `json:"userId"`
	Page   int    `json:"page"`
	Limit  int    `json:"limit"`
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
	apps, err := persistence.ListAppsPaginated(ctx, s.db, req.UserID, req.Page, req.Limit)
	if err != nil {
		return nil, err
	}

	return &ListAppsResp{Apps: apps}, nil
}
