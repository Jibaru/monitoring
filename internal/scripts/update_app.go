package scripts

import (
	"context"

	"monitoring/internal/domain"
)

type UpdateAppReq struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	AppKey string `json:"appKey"`
}

type UpdateAppResp struct {
	domain.App
}

type UpdateAppScript struct {
	appRepo domain.AppRepo
}

func NewUpdateAppScript(appRepo domain.AppRepo) *UpdateAppScript {
	return &UpdateAppScript{appRepo: appRepo}
}

func (s *UpdateAppScript) Exec(ctx context.Context, req UpdateAppReq) (*UpdateAppResp, error) {
	id, err := domain.NewID(req.ID)
	if err != nil {
		return nil, err
	}

	app, err := s.appRepo.GetAppByID(ctx, id)
	if err != nil {
		return nil, err
	}

	err = app.ChangeName(req.Name)
	if err != nil {
		return nil, err
	}

	err = app.ChangeAppKey(req.AppKey)
	if err != nil {
		return nil, err
	}

	err = s.appRepo.UpdateApp(ctx, *app)
	if err != nil {
		return nil, err
	}

	return &UpdateAppResp{
		App: *app,
	}, nil
}
