package scripts

import (
	"context"
	"fmt"
	"monitoring/internal/domain"
	"time"

	"github.com/google/uuid"
)

type VisitorLoginResp struct {
	LoginResp
}

type VisitorLoginScript struct {
	userRepo  domain.UserRepo
	jwtSecret []byte
}

func NewVisitorLoginScript(userRepo domain.UserRepo, jwtSecret string) *VisitorLoginScript {
	return &VisitorLoginScript{
		userRepo:  userRepo,
		jwtSecret: []byte(jwtSecret),
	}
}

func (s *VisitorLoginScript) Exec(ctx context.Context) (*VisitorLoginResp, error) {
	visitorEmail := uuid.NewString() + "_monitoring_" + fmt.Sprintf("%v", Now().Unix()) + "@mail.app"

	user, err := domain.NewUser(
		domain.NewAutoID(),
		generateUsername(),
		visitorEmail,
		"",
		Now().UTC(),
		"",
		time.Time{},
		nil,
		true,
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}

	err = s.userRepo.SaveUser(ctx, *user)
	if err != nil {
		return nil, err
	}

	tokenString, err := generateToken(user.ID(), user.Email(), s.jwtSecret)
	if err != nil {
		return nil, err
	}

	return &VisitorLoginResp{
		LoginResp: *userToLoginResp(tokenString, user),
	}, nil
}
