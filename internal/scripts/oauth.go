package scripts

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/mongo"

	"monitoring/internal/domain"
)

type OAuthReq struct {
	Username string `json:"username"`
	Email    string `json:"email"`
}

type OAuthResp struct {
	LoginResp
}

type OAuthScript struct {
	userRepo  domain.UserRepo
	jwtSecret []byte
}

func NewOAuthScript(
	userRepo domain.UserRepo,
	jwtSecret []byte,
) *OAuthScript {
	return &OAuthScript{userRepo: userRepo, jwtSecret: jwtSecret}
}

func (s *OAuthScript) Exec(ctx context.Context, req OAuthReq) (*OAuthResp, error) {
	user, err := s.userRepo.GetUserByEmail(ctx, req.Email)
	if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
		return nil, err
	}

	if err != nil && errors.Is(err, mongo.ErrNoDocuments) {
		// Register as root
		validatedAt := Now().UTC()
		newUser, err := domain.NewUser(
			domain.NewAutoID(),
			req.Username,
			req.Email,
			"",
			Now().UTC(),
			"",
			time.Time{},
			&validatedAt,
			false,
			true,
			nil,
		)
		if err != nil {
			return nil, err
		}

		err = s.userRepo.SaveUser(ctx, *newUser)
		if err != nil {
			return nil, err
		}

		user = newUser
	}

	tokenString, err := generateToken(user.ID(), user.Email(), s.jwtSecret)
	if err != nil {
		return nil, err
	}

	return &OAuthResp{
		LoginResp: *userToLoginResp(tokenString, user),
	}, nil
}
