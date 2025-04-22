package scripts

import (
	"context"
	"errors"

	"monitoring/internal/domain"
)

type DeleteUserReq struct {
	RootUserID string `json:"rootUserId"`
	UserID     string `json:"userId"`
}

type DeleteUserScript struct {
	userRepo domain.UserRepo
}

func NewDeleteUserScript(userRepo domain.UserRepo) *DeleteUserScript {
	return &DeleteUserScript{userRepo: userRepo}
}

func (s *DeleteUserScript) Exec(ctx context.Context, req DeleteUserReq) error {
	id, err := domain.NewID(req.UserID)
	if err != nil {
		return err
	}

	rootUserID, err := domain.NewID(req.RootUserID)
	if err != nil {
		return err
	}

	user, err := s.userRepo.GetUserByID(ctx, id)
	if err != nil {
		return err
	}

	if user.RootUserID() != nil && *user.RootUserID() != rootUserID {
		return errors.New("user does not belong to the root user")
	}

	if user.RootUserID() == nil {
		return errors.New("cannot delete a root user")
	}

	err = s.userRepo.DeleteUser(ctx, id)
	if err != nil {
		return err
	}
	return nil
}
