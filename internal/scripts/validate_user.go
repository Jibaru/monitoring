package scripts

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"monitoring/internal/persistence"
)

var (
	ErrValidateUserScriptValidationExpired = errors.New("validation expired")
	ErrValidateUserScriptInvalidPin        = errors.New("invalid pin")
)

type ValidateUserReq struct {
	UserID string `-:"userId"`
	Pin    string `json:"pin"`
}

type ValidateUserResp struct {
	ValidatedAt string `json:"validatedAt"`
}

type ValidateUserScript struct {
	db *mongo.Database
}

func NewValidateUserScript(
	db *mongo.Database,
) *ValidateUserScript {
	return &ValidateUserScript{db: db}
}

func (s *ValidateUserScript) Exec(ctx context.Context, req ValidateUserReq) (*ValidateUserResp, error) {
	userID, err := primitive.ObjectIDFromHex(req.UserID)
	if err != nil {
		return nil, err
	}

	user, err := persistence.GetUserByID(ctx, s.db, userID)
	if err != nil {
		return nil, err
	}

	if user.PinExpiresAt.Before(time.Now().UTC()) {
		return nil, ErrValidateUserScriptValidationExpired
	}

	if user.Pin != req.Pin {
		return nil, ErrValidateUserScriptInvalidPin
	}

	validatedAt := time.Now().UTC()
	user.ValidatedAt = &validatedAt

	err = persistence.UpdateUser(ctx, s.db, *user)
	if err != nil {
		return nil, err
	}

	return &ValidateUserResp{
		ValidatedAt: validatedAt.Format(time.RFC3339),
	}, nil
}
