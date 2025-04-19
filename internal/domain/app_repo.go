package domain

import (
	"context"
)

type AppRepo interface {
	SaveApp(ctx context.Context, app App) error
	UpdateApp(ctx context.Context, app App) error
	GetAppByID(ctx context.Context, appID ID) (*App, error)
	GetAppByKey(ctx context.Context, appKey string) (*App, error)
	DeleteApp(ctx context.Context, appID ID) error
	ListApps(ctx context.Context, criteria Criteria) ([]App, error)
}
