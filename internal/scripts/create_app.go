package scripts

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/mongo"

	"monitoring/internal/domain"
)

type CreateAppReq struct {
	Name   string `json:"name"`
	AppKey string `json:"appKey"`
	UserID string `json:"userId"`
}

type CreateAppResp struct {
	domain.App
}

type CreateAppScript struct {
	appRepo domain.AppRepo
}

func NewCreateAppScript(appRepo domain.AppRepo) *CreateAppScript {
	return &CreateAppScript{appRepo: appRepo}
}

func (s *CreateAppScript) Exec(ctx context.Context, req CreateAppReq) (*CreateAppResp, error) {
	existing, err := s.appRepo.GetAppByKey(ctx, req.AppKey)
	if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
		return nil, err
	}

	if existing != nil {
		return nil, errors.New("app with the provided app key already exists")
	}

	userID, err := domain.NewID(req.UserID)
	if err != nil {
		return nil, err
	}

	app, err := domain.NewApp(
		domain.NewAutoID(),
		req.Name,
		req.AppKey,
		userID,
		Now().UTC(),
	)
	if err != nil {
		return nil, err
	}

	err = s.appRepo.SaveApp(ctx, *app)
	if err != nil {
		return nil, err
	}

	return &CreateAppResp{
		App: *app,
	}, nil
}
