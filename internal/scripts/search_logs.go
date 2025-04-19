package scripts

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/openai/openai-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"monitoring/internal/domain"
	"monitoring/internal/persistence"
)

type SearchLogsReq struct {
	UserID     string    `json:"-"`
	Page       int       `form:"page"`
	Limit      int       `form:"limit"`
	SortOrder  string    `form:"sortOrder"`
	SearchTerm string    `form:"searchTerm"`
	LogLevel   string    `form:"logLevel"`
	From       time.Time `form:"from"`
	To         time.Time `form:"to"`
	AppID      string    `form:"appId"`
}

type SearchLogsResp struct {
	Data []persistence.Log `json:"data"`
}

type SearchLogsScript struct {
	db           *mongo.Database
	openaiClient *openai.Client
	appRepo      domain.AppRepo
}

func NewSearchLogsScript(db *mongo.Database, appRepo domain.AppRepo) *SearchLogsScript {
	// TODO: add pipeline generation using openai client
	client := openai.NewClient(openai.DefaultClientOptions()...)
	return &SearchLogsScript{db: db, openaiClient: &client, appRepo: appRepo}
}

func (s *SearchLogsScript) Exec(ctx context.Context, req SearchLogsReq) (*SearchLogsResp, error) {
	userID, err := domain.NewID(req.UserID)
	if err != nil {
		return nil, err
	}

	apps, err := s.appRepo.ListApps(ctx, domain.NewCriteria(
		[]domain.Filter{
			domain.NewFilter("userId", domain.Equals, userID),
		},
		domain.EmptyPagination,
		domain.EmptySort,
	))
	if err != nil {
		return nil, err
	}

	filters := []persistence.Filter{}

	if strings.TrimSpace(req.SearchTerm) != "" {
		filters = append(filters, persistence.NewFilter("raw", persistence.Like, req.SearchTerm))
	}

	if strings.TrimSpace(req.LogLevel) != "" {
		filters = append(filters, persistence.NewFilter("level", persistence.Equals, req.LogLevel))
	}

	if !req.From.IsZero() {
		filters = append(filters, persistence.NewFilter("timestamp", persistence.GreaterThanOrEqual, req.From.UTC()))
	}

	if !req.To.IsZero() {
		filters = append(filters, persistence.NewFilter("timestamp", persistence.LessThanOrEqual, req.To.UTC()))
	}

	appsIDs := make(bson.A, len(apps))
	for i, app := range apps {
		appsIDs[i] = app.ID
	}
	if strings.TrimSpace(req.AppID) != "" {
		appID, err := primitive.ObjectIDFromHex(req.AppID)
		if err != nil {
			return nil, err
		}

		exists := false
		for _, existingAppID := range appsIDs {
			if existingAppID == appID {
				exists = true
				break
			}
		}

		if !exists {
			return nil, fmt.Errorf("app with ID %s does not exist for the user", req.AppID)
		}

		appsIDs = bson.A{appID}
	}

	filters = append(filters, persistence.NewFilter("appId", persistence.In, appsIDs))

	criteria := persistence.NewCriteria(
		filters,
		persistence.NewPagination(req.Limit, (req.Page-1)*req.Limit),
		persistence.NewSort("timestamp", persistence.SortOrder(req.SortOrder)),
	)

	logs, err := persistence.ListLogs(ctx, s.db, criteria)
	if err != nil {
		return nil, err
	}

	return &SearchLogsResp{Data: logs}, nil
}
