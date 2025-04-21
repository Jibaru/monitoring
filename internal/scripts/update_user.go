package scripts

import (
	"context"

	"monitoring/internal/domain"
)

type UpdateUserReq struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}

type UpdateUserResp struct {
	Username string `json:"username"`
}

type UpdateUserScript struct {
	userRepo domain.UserRepo
}

func NewUpdateUserScript(userRepo domain.UserRepo) *UpdateUserScript {
	return &UpdateUserScript{userRepo: userRepo}
}

func (s *UpdateUserScript) Exec(ctx context.Context, req UpdateUserReq) (*UpdateUserResp, error) {
	id, err := domain.NewID(req.ID)
	if err != nil {
		return nil, err
	}

	user, err := s.userRepo.GetUserByID(ctx, id)
	if err != nil {
		return nil, err
	}

	err = user.ChangeUsername(req.Username)
	if err != nil {
		return nil, err
	}

	err = s.userRepo.UpdateUser(ctx, *user)
	if err != nil {
		return nil, err
	}

	return &UpdateUserResp{
		Username: user.Username(),
	}, nil
}
