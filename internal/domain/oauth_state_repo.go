package domain

import (
	"context"
	"errors"
)

var (
	ErrNoOAuthStatesDeleted = errors.New("no oauth states deleted")
)

type OAuthStateRepo interface {
	SaveOAuthState(ctx context.Context, oauthState OAuthState) error
	DeleteOAuthStateByState(ctx context.Context, state string) error
}
