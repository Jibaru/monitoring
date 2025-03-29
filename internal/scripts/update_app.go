package scripts

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"monitoring/internal/persistence"
)

type UpdateAppReq struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	AppKey string `json:"appKey"`
}

type UpdateAppResp struct {
	persistence.App
}

type UpdateAppScript struct {
	db *mongo.Database
}

func NewUpdateAppScript(db *mongo.Database) *UpdateAppScript {
	return &UpdateAppScript{db: db}
}

func (s *UpdateAppScript) Exec(ctx context.Context, req UpdateAppReq) (*UpdateAppResp, error) {
	id, err := primitive.ObjectIDFromHex(req.ID)
	if err != nil {
		return nil, err
	}

	app, err := persistence.GetAppByID(ctx, s.db, id)
	if err != nil {
		return nil, err
	}

	app.Name = req.Name
	app.AppKey = req.AppKey

	err = persistence.UpdateApp(ctx, s.db, *app)
	if err != nil {
		return nil, err
	}

	return &UpdateAppResp{
		App: *app,
	}, nil
}
