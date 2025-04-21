package domain

import (
	"context"
)

type LogSchemaRepo interface {
	Get(ctx context.Context, userID ID, appIDs []ID, optionalRange *Range) (LogSchema, error)
}
