package scripts

import (
	"context"
	"errors"
	"time"

	"monitoring/internal/domain"
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
	userRepo domain.UserRepo
}

func NewValidateUserScript(
	userRepo domain.UserRepo,
) *ValidateUserScript {
	return &ValidateUserScript{userRepo: userRepo}
}

func (s *ValidateUserScript) Exec(ctx context.Context, req ValidateUserReq) (*ValidateUserResp, error) {
	userID, err := domain.NewID(req.UserID)
	if err != nil {
		return nil, err
	}

	user, err := s.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if user.PinExpiresAt().Before(Now().UTC()) {
		return nil, ErrValidateUserScriptValidationExpired
	}

	if user.Pin() != req.Pin {
		return nil, ErrValidateUserScriptInvalidPin
	}

	validatedAt := Now().UTC()
	err = user.ChangeValidatedAt(&validatedAt)
	if err != nil {
		return nil, err
	}

	err = s.userRepo.UpdateUser(ctx, *user)
	if err != nil {
		return nil, err
	}

	return &ValidateUserResp{
		ValidatedAt: validatedAt.Format(time.RFC3339),
	}, nil
}
