package scripts

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/oauth2"

	"monitoring/config"
	"monitoring/internal/domain"
	"monitoring/internal/domain/services"
)

type FinishOAuthReq struct {
	State string `json:"state"`
	Code  string `json:"code"`
}

type FinishOAuthResp struct {
	URL string
}

type FinishOAuthScript struct {
	userRepo           domain.UserRepo
	oauthStateRepo     domain.OAuthStateRepo
	oauthCfg           *oauth2.Config
	oauthInfoExtractor services.OAuthInfoExtractor
	cfg                config.Config
}

func NewFinishOAuthScript(
	userRepo domain.UserRepo,
	oauthStateRepo domain.OAuthStateRepo,
	oauthCfg *oauth2.Config,
	oauthInfoExtractor services.OAuthInfoExtractor,
	cfg config.Config,
) *FinishOAuthScript {
	return &FinishOAuthScript{
		userRepo:           userRepo,
		oauthStateRepo:     oauthStateRepo,
		oauthCfg:           oauthCfg,
		oauthInfoExtractor: oauthInfoExtractor,
		cfg:                cfg,
	}
}

func (s *FinishOAuthScript) Exec(ctx context.Context, req FinishOAuthReq) (*FinishOAuthResp, error) {
	err := s.oauthStateRepo.DeleteOAuthStateByState(ctx, req.State)
	if err != nil && !errors.Is(err, domain.ErrNoOAuthStatesDeleted) {
		return nil, err
	}

	if err != nil && errors.Is(err, domain.ErrNoOAuthStatesDeleted) {
		return nil, errors.New("invalid state")
	}

	token, err := s.oauthCfg.Exchange(ctx, req.Code)
	if err != nil {
		return nil, err
	}

	username, email, err := s.oauthInfoExtractor(token.AccessToken)
	if err != nil {
		return nil, err
	}

	user, err := s.userRepo.GetUserByEmail(ctx, email)
	if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
		return nil, err
	}

	if err != nil && errors.Is(err, mongo.ErrNoDocuments) {
		// Register as root
		validatedAt := Now().UTC()
		newUser, err := domain.NewUser(
			domain.NewAutoID(),
			username,
			email,
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

	tokenString, err := generateToken(user.ID(), user.Email(), []byte(s.cfg.JWTSecret))
	if err != nil {
		return nil, err
	}

	isVisitor := "false"
	if user.IsVisitor() {
		isVisitor = "true"
	}

	url := fmt.Sprintf("%s/login?token=%s&id=%s&email=%s&username=%s&isVisitor=%s",
		s.cfg.WebBaseURI,
		tokenString,
		user.ID(),
		user.Email(),
		user.Username(),
		isVisitor,
	)

	return &FinishOAuthResp{
		URL: url,
	}, nil
}
