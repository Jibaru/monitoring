package scripts

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"monitoring/internal/persistence"
)

type CreateAppReq struct {
	Name   string `json:"name"`
	AppKey string `json:"appKey"`
	UserID string `json:"userId"`
}

type CreateAppResp struct {
	persistence.App
}

type CreateAppScript struct {
	db *mongo.Database
}

func NewCreateAppScript(db *mongo.Database) *CreateAppScript {
	return &CreateAppScript{db: db}
}

func (s *CreateAppScript) Exec(ctx context.Context, req CreateAppReq) (*CreateAppResp, error) {
	existing, err := persistence.GetAppByKey(ctx, s.db, req.AppKey)
	if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
		return nil, err
	}

	if existing != nil {
		return nil, errors.New("app with the provided app key already exists")
	}

	app := persistence.App{
		ID:        primitive.NewObjectID(),
		AppKey:    req.AppKey,
		Name:      req.Name,
		UserID:    req.UserID,
		CreatedAt: time.Now().UTC(),
	}
	err = persistence.SaveApp(ctx, s.db, app)
	if err != nil {
		return nil, err
	}

	return &CreateAppResp{
		App: app,
	}, nil
}
