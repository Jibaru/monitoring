package domain

import (
	"context"
)

type LogSchemaResult struct {
	Total  int            `json:"total"`
	Schema map[string]int `json:"schema"`
}

type LogSchemaRepo interface {
	GetLogsSchema(ctx context.Context, userID ID, appIDs []ID, optionalRange *Range) (LogSchemaResult, error)
}
