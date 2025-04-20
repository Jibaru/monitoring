package scripts

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/openai/openai-go"

	"monitoring/internal/domain"
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
	Data []domain.Log `json:"data"`
}

type SearchLogsScript struct {
	openaiClient *openai.Client
	appRepo      domain.AppRepo
	logRepo      domain.LogRepo
}

func NewSearchLogsScript(appRepo domain.AppRepo, logRepo domain.LogRepo) *SearchLogsScript {
	// TODO: add pipeline generation using openai client
	client := openai.NewClient(openai.DefaultClientOptions()...)
	return &SearchLogsScript{openaiClient: &client, appRepo: appRepo, logRepo: logRepo}
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

	filters := []domain.Filter{}

	if strings.TrimSpace(req.SearchTerm) != "" {
		filters = append(filters, domain.NewFilter("raw", domain.Like, req.SearchTerm))
	}

	if strings.TrimSpace(req.LogLevel) != "" {
		filters = append(filters, domain.NewFilter("level", domain.Equals, req.LogLevel))
	}

	if !req.From.IsZero() {
		filters = append(filters, domain.NewFilter("timestamp", domain.GreaterThanOrEqual, req.From.UTC()))
	}

	if !req.To.IsZero() {
		filters = append(filters, domain.NewFilter("timestamp", domain.LessThanOrEqual, req.To.UTC()))
	}

	appsIDs := make([]any, len(apps))
	for i, app := range apps {
		appsIDs[i] = app.ID()
	}
	if strings.TrimSpace(req.AppID) != "" {
		appID, err := domain.NewID(req.AppID)
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

		appsIDs = []any{appID}
	}

	filters = append(filters, domain.NewFilter("appId", domain.In, appsIDs))

	criteria := domain.NewCriteria(
		filters,
		domain.NewPagination(req.Limit, (req.Page-1)*req.Limit),
		domain.NewSort("timestamp", domain.SortOrder(req.SortOrder)),
	)

	logs, err := s.logRepo.ListLogs(ctx, criteria)
	if err != nil {
		return nil, err
	}

	return &SearchLogsResp{Data: logs}, nil
}
