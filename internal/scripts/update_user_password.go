package scripts

import (
	"context"
	"errors"

	"monitoring/internal/domain"
)

type UpdateUserPasswordReq struct {
	ID          string `json:"id"`
	OldPassword string `json:"oldPassword"`
	NewPassword string `json:"newPassword"`
}

type UpdateUserPasswordScript struct {
	userRepo domain.UserRepo
}

func NewUpdateUserPasswordScript(userRepo domain.UserRepo) *UpdateUserPasswordScript {
	return &UpdateUserPasswordScript{userRepo: userRepo}
}

func (s *UpdateUserPasswordScript) Exec(ctx context.Context, req UpdateUserPasswordReq) error {
	id, err := domain.NewID(req.ID)
	if err != nil {
		return err
	}

	user, err := s.userRepo.GetUserByID(ctx, id)
	if err != nil {
		return err
	}

	if user.Password() != "" && !isValidPassword(user.Password(), req.OldPassword) {
		return errors.New("old password is invalid")
	}

	encryptedPassword, err := encryptPassword(req.NewPassword)
	if err != nil {
		return err
	}

	err = user.ChangePassword(encryptedPassword)
	if err != nil {
		return err
	}

	err = s.userRepo.UpdateUser(ctx, *user)
	if err != nil {
		return err
	}

	return nil
}
