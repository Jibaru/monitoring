package scripts

import (
	"context"
	"crypto/rand"
	"encoding/hex"

	"golang.org/x/oauth2"

	"monitoring/internal/domain"
)

type StartOAuthScript struct {
	oauthStateRepo domain.OAuthStateRepo
	cfg            *oauth2.Config
}

type StartOAuthResp struct {
	URL string
}

func NewStartOAuthScript(oauthStateRepo domain.OAuthStateRepo, cfg *oauth2.Config) *StartOAuthScript {
	return &StartOAuthScript{oauthStateRepo: oauthStateRepo, cfg: cfg}
}

func (s *StartOAuthScript) Exec(ctx context.Context) (*StartOAuthResp, error) {
	state, err := s.generateState(16)
	if err != nil {
		return nil, err
	}

	oauthState, err := domain.NewOAuthState(domain.NewAutoID(), state)
	if err != nil {
		return nil, err
	}

	err = s.oauthStateRepo.SaveOAuthState(ctx, *oauthState)
	if err != nil {
		return nil, err
	}

	url := s.cfg.AuthCodeURL(state)

	return &StartOAuthResp{URL: url}, nil
}

func (s *StartOAuthScript) generateState(length int) (string, error) {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
