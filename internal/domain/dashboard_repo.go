package domain

import (
	"context"
	"time"
)

type DashboardOverviewKPI struct {
	Total      int64   `bson:"total"`
	Percentage float64 `bson:"percentage"`
}

type Period struct {
	Year  int `bson:"year"`
	Month int `bson:"month"`
}

type DashboardOverviewKPIs struct {
	Logs       DashboardOverviewKPI `bson:"logs"`
	Errors     DashboardOverviewKPI `bson:"errors"`
	Warnings   DashboardOverviewKPI `bson:"warnings"`
	Info       DashboardOverviewKPI `bson:"info"`
	LogsPerApp map[string]struct {
		AppName string `bson:"appName"`
		Total   int64  `bson:"total"`
	} `bson:"logsPerApp"`
	LogsByPeriod map[Period]struct {
		Total int64 `bson:"total"`
	} `bson:"logsByPeriod"`
}

type Range struct {
	From time.Time
	To   time.Time
}

type DashboardRepo interface {
	GetDashboardOverviewKPIs(
		ctx context.Context,
		userID ID,
		optionalRange *Range,
	) (DashboardOverviewKPIs, error)
}
