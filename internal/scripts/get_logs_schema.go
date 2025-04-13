package scripts

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"monitoring/internal/persistence"
)

type GetLogsSchemaReq struct {
	UserID string `json:"-"`
}

type GetLogsSchemaResp struct {
	persistence.LogSchemaResult
}

type GetLogsSchemaScript struct {
	db *mongo.Database
}

func NewGetLogsSchemaScript(db *mongo.Database) *GetLogsSchemaScript {
	return &GetLogsSchemaScript{db: db}
}

func (s *GetLogsSchemaScript) Exec(ctx context.Context, req GetLogsSchemaReq) (*GetLogsSchemaResp, error) {
	userID, err := primitive.ObjectIDFromHex(req.UserID)
	if err != nil {
		return nil, err
	}

	// TODO: add appIDs filter and [from, to] timestamp filter
	result, err := persistence.GetLogsSchema(ctx, s.db, userID, nil, nil)
	if err != nil {
		return nil, err
	}

	return &GetLogsSchemaResp{
		LogSchemaResult: result,
	}, nil
}
