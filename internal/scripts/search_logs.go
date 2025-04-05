package scripts

import (
	"context"
	"strings"

	"github.com/openai/openai-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"monitoring/internal/persistence"
)

type SearchLogsReq struct {
	UserID     string `json:"-"`
	Page       int    `form:"page"`
	Limit      int    `form:"limit"`
	SortOrder  string `form:"sortOrder"`
	SearchTerm string `form:"searchTerm"`
	LogLevel   string `form:"logLevel"`
}

type SearchLogsResp struct {
	Data []persistence.Log `json:"data"`
}

type SearchLogsScript struct {
	db           *mongo.Database
	openaiClient *openai.Client
}

func NewSearchLogsScript(db *mongo.Database) *SearchLogsScript {
	// TODO: add pipeline generation using openai client
	client := openai.NewClient(openai.DefaultClientOptions()...)
	return &SearchLogsScript{db: db, openaiClient: &client}
}

func (s *SearchLogsScript) Exec(ctx context.Context, req SearchLogsReq) (*SearchLogsResp, error) {
	userID, err := primitive.ObjectIDFromHex(req.UserID)
	if err != nil {
		return nil, err
	}

	apps, err := persistence.ListApps(ctx, s.db, persistence.NewCriteria(
		[]persistence.Filter{
			persistence.NewFilter("userId", persistence.Equals, userID),
		},
		persistence.EmptyPagination,
		persistence.EmptySort,
	))
	if err != nil {
		return nil, err
	}

	appsIDs := make(bson.A, len(apps))
	for i, app := range apps {
		appsIDs[i] = app.ID
	}

	filters := []persistence.Filter{
		persistence.NewFilter("appId", persistence.In, appsIDs),
	}

	if strings.TrimSpace(req.SearchTerm) != "" {
		filters = append(filters, persistence.NewFilter("raw", persistence.Like, req.SearchTerm))
	}

	if strings.TrimSpace(req.LogLevel) != "" {
		filters = append(filters, persistence.NewFilter("level", persistence.Equals, req.LogLevel))
	}

	criteria := persistence.NewCriteria(
		filters,
		persistence.NewPagination(req.Limit, (req.Page-1)*req.Limit),
		persistence.NewSort("createdAt", persistence.SortOrder(req.SortOrder)),
	)

	logs, err := persistence.ListLogs(ctx, s.db, criteria)
	if err != nil {
		return nil, err
	}

	return &SearchLogsResp{Data: logs}, nil
}
