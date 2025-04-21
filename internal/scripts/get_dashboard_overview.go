package scripts

import (
	"context"
	"fmt"
	"time"

	"monitoring/internal/domain"
)

type GetDashboardOverviewReq struct {
	UserID string    `json:"-"`
	From   time.Time `form:"from"`
	To     time.Time `form:"to"`
}

type KPI struct {
	Total      int64   `json:"total"`
	Percentage float64 `json:"percentage"`
}

type GetDashboardOverviewResp struct {
	Logs       KPI `json:"logs"`
	Errors     KPI `json:"errors"`
	Warnings   KPI `json:"warnings"`
	Info       KPI `json:"info"`
	LogsPerApp map[string]struct {
		AppName string `json:"appName"`
		Total   int64  `json:"total"`
	} `json:"logsPerApp"`
	LogsByPeriod map[string]struct {
		Total int64 `json:"total"`
	} `json:"logsByPeriod"`
}

type GetDashboardOverviewScript struct {
	dashboardRepo domain.DashboardRepo
}

func NewGetDashboardOverviewScript(dashboardRepo domain.DashboardRepo) *GetDashboardOverviewScript {
	return &GetDashboardOverviewScript{dashboardRepo: dashboardRepo}
}

func (s *GetDashboardOverviewScript) Exec(ctx context.Context, req GetDashboardOverviewReq) (*GetDashboardOverviewResp, error) {
	userID, err := domain.NewID(req.UserID)
	if err != nil {
		return nil, err
	}

	var dateRange *domain.Range
	if !req.From.IsZero() || !req.To.IsZero() {
		dateRange = &domain.Range{}
		if !req.From.IsZero() {
			dateRange.From = req.From.UTC()
		}

		if !req.To.IsZero() {
			dateRange.To = req.To.UTC()
		}
	}

	kpis, err := s.dashboardRepo.OverviewKPIs(ctx, userID, dateRange)
	if err != nil {
		return nil, err
	}

	resp := &GetDashboardOverviewResp{
		Logs: KPI{
			Total:      kpis.Logs.Total,
			Percentage: kpis.Logs.Percentage,
		},
		Errors: KPI{
			Total:      kpis.Errors.Total,
			Percentage: kpis.Errors.Percentage,
		},
		Warnings: KPI{
			Total:      kpis.Warnings.Total,
			Percentage: kpis.Warnings.Percentage,
		},
		Info: KPI{
			Total:      kpis.Info.Total,
			Percentage: kpis.Info.Percentage,
		},
		LogsPerApp: make(map[string]struct {
			AppName string `json:"appName"`
			Total   int64  `json:"total"`
		}, len(kpis.LogsPerApp)),
		LogsByPeriod: make(map[string]struct {
			Total int64 `json:"total"`
		}, len(kpis.LogsByPeriod)),
	}

	for appID, app := range kpis.LogsPerApp {
		resp.LogsPerApp[appID] = struct {
			AppName string `json:"appName"`
			Total   int64  `json:"total"`
		}{
			AppName: app.AppName,
			Total:   app.Total,
		}
	}

	for period, periodData := range kpis.LogsByPeriod {
		periodStr := fmt.Sprintf("%d-%02d", period.Year, period.Month)
		resp.LogsByPeriod[periodStr] = struct {
			Total int64 `json:"total"`
		}{
			Total: periodData.Total,
		}
	}

	return resp, nil
}
