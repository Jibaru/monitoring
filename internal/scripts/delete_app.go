package scripts

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"monitoring/internal/persistence"
)

type DeleteAppReq struct {
	AppID string `json:"app_id"`
}

type DeleteAppScript struct {
	db *mongo.Database
}

func NewDeleteAppScript(db *mongo.Database) *DeleteAppScript {
	return &DeleteAppScript{db: db}
}

func (s *DeleteAppScript) Exec(ctx context.Context, req DeleteAppReq) error {
	id, err := primitive.ObjectIDFromHex(req.AppID)
	if err != nil {
		return err
	}

	err = persistence.DeleteApp(ctx, s.db, id)
	if err != nil {
		return err
	}
	return nil
}
