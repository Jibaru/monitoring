package scripts

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"monitoring/internal/persistence"
)

type UpdateUserPasswordReq struct {
	ID          string `json:"id"`
	OldPassword string `json:"oldPassword"`
	NewPassword string `json:"newPassword"`
}

type UpdateUserPasswordScript struct {
	db *mongo.Database
}

func NewUpdateUserPasswordScript(db *mongo.Database) *UpdateUserPasswordScript {
	return &UpdateUserPasswordScript{db: db}
}

func (s *UpdateUserPasswordScript) Exec(ctx context.Context, req UpdateUserPasswordReq) error {
	id, err := primitive.ObjectIDFromHex(req.ID)
	if err != nil {
		return err
	}

	user, err := persistence.GetUserByID(ctx, s.db, id)
	if err != nil {
		return err
	}

	if user.Password != "" && !isValidPassword(user.Password, req.OldPassword) {
		return errors.New("old password is invalid")
	}

	encryptedPassword, err := encryptPassword(req.NewPassword)
	if err != nil {
		return err
	}

	user.Password = encryptedPassword

	err = persistence.UpdateUser(ctx, s.db, *user)
	if err != nil {
		return err
	}

	return nil
}
