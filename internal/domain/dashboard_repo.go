package domain

import (
	"context"
	"time"
)

type Range struct {
	From time.Time
	To   time.Time
}

type DashboardRepo interface {
	OverviewKPIs(
		ctx context.Context,
		userID ID,
		optionalRange *Range,
	) (DashboardOverviewKPIs, error)
}
