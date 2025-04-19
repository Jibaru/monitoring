package domain

import (
	"context"
)

type LogRepo interface {
	SaveLogs(ctx context.Context, logs []Log) error
	ListLogs(ctx context.Context, criteria Criteria) ([]Log, error)
}
