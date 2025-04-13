package scripts

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"monitoring/internal/persistence"
)

type UpdateUserReq struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}

type UpdateUserResp struct {
	Username string `json:"username"`
}

type UpdateUserScript struct {
	db *mongo.Database
}

func NewUpdateUserScript(db *mongo.Database) *UpdateUserScript {
	return &UpdateUserScript{db: db}
}

func (s *UpdateUserScript) Exec(ctx context.Context, req UpdateUserReq) (*UpdateUserResp, error) {
	id, err := primitive.ObjectIDFromHex(req.ID)
	if err != nil {
		return nil, err
	}

	user, err := persistence.GetUserByID(ctx, s.db, id)
	if err != nil {
		return nil, err
	}

	user.Username = req.Username

	err = persistence.UpdateUser(ctx, s.db, *user)
	if err != nil {
		return nil, err
	}

	return &UpdateUserResp{
		Username: user.Username,
	}, nil
}
