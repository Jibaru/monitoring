package scripts

import (
	"context"
	"errors"
	"monitoring/internal/domain"
)

type LoginReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResp struct {
	Token string `json:"token"`
	User  struct {
		ID        string `json:"id"`
		Username  string `json:"username"`
		Email     string `json:"email"`
		IsVisitor bool   `json:"isVisitor"`
		FromOAuth bool   `json:"fromOAuth"`
	} `json:"user"`
}

type LoginScript struct {
	userRepo  domain.UserRepo
	jwtSecret []byte
}

func NewLoginScript(userRepo domain.UserRepo, jwtSecret string) *LoginScript {
	return &LoginScript{
		userRepo:  userRepo,
		jwtSecret: []byte(jwtSecret),
	}
}

func (s *LoginScript) Exec(ctx context.Context, req LoginReq) (*LoginResp, error) {
	user, err := s.userRepo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}

	if user.ValidatedAt() == nil {
		return nil, errors.New("user should be validated to login")
	}

	if !isValidPassword(user.Password(), req.Password) {
		return nil, errors.New("email or password are invalid")
	}

	tokenString, err := generateToken(user.ID(), user.Email(), s.jwtSecret)
	if err != nil {
		return nil, err
	}

	return userToLoginResp(tokenString, user), nil
}

func userToLoginResp(token string, user *domain.User) *LoginResp {
	return &LoginResp{
		Token: token,
		User: struct {
			ID        string `json:"id"`
			Username  string `json:"username"`
			Email     string `json:"email"`
			IsVisitor bool   `json:"isVisitor"`
			FromOAuth bool   `json:"fromOAuth"`
		}{
			ID:        user.ID().Hex(),
			Username:  user.Username(),
			Email:     user.Email(),
			IsVisitor: user.IsVisitor(),
			FromOAuth: user.FromOAuth(),
		},
	}
}
