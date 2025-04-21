package scripts

import (
	"context"

	"monitoring/internal/domain"
)

type GetLogsSchemaReq struct {
	UserID string `json:"-"`
}

type GetLogsSchemaResp struct {
	domain.LogSchema
}

type GetLogsSchemaScript struct {
	logSchemaRepo domain.LogSchemaRepo
}

func NewGetLogsSchemaScript(logSchemaRepo domain.LogSchemaRepo) *GetLogsSchemaScript {
	return &GetLogsSchemaScript{logSchemaRepo: logSchemaRepo}
}

func (s *GetLogsSchemaScript) Exec(ctx context.Context, req GetLogsSchemaReq) (*GetLogsSchemaResp, error) {
	userID, err := domain.NewID(req.UserID)
	if err != nil {
		return nil, err
	}

	// TODO: add appIDs filter and [from, to] timestamp filter
	result, err := s.logSchemaRepo.Get(ctx, userID, nil, nil)
	if err != nil {
		return nil, err
	}

	return &GetLogsSchemaResp{
		LogSchema: result,
	}, nil
}
