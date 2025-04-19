package scripts

import (
	"context"

	"monitoring/internal/domain"
)

type DeleteAppReq struct {
	AppID string `json:"app_id"`
}

type DeleteAppScript struct {
	appRepo domain.AppRepo
}

func NewDeleteAppScript(appRepo domain.AppRepo) *DeleteAppScript {
	return &DeleteAppScript{appRepo: appRepo}
}

func (s *DeleteAppScript) Exec(ctx context.Context, req DeleteAppReq) error {
	id, err := domain.NewID(req.AppID)
	if err != nil {
		return err
	}

	err = s.appRepo.DeleteApp(ctx, id)
	if err != nil {
		return err
	}
	return nil
}
