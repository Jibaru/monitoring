package scripts

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"monitoring/internal/persistence"
)

type OAuthReq struct {
	Username string `json:"username"`
	Email    string `json:"email"`
}

type OAuthResp struct {
	LoginResp
}

type OAuthScript struct {
	db        *mongo.Database
	jwtSecret []byte
}

func NewOAuthScript(
	db *mongo.Database,
	jwtSecret []byte,
) *OAuthScript {
	return &OAuthScript{db: db, jwtSecret: jwtSecret}
}

func (s *OAuthScript) Exec(ctx context.Context, req OAuthReq) (*OAuthResp, error) {
	user, err := persistence.GetUserByEmail(ctx, s.db, req.Email)
	if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
		return nil, err
	}

	if err != nil && errors.Is(err, mongo.ErrNoDocuments) {
		// Register
		validatedAt := time.Now().UTC()
		u := persistence.User{
			ID:           primitive.NewObjectID(),
			Email:        req.Email,
			Password:     "",
			RegisteredAt: time.Now().UTC(),
			ValidatedAt:  &validatedAt,
			IsVisitor:    false,
			FromGithub:   true,
		}

		err = persistence.SaveUser(ctx, s.db, u)
		if err != nil {
			return nil, err
		}

		user = &u
	}

	tokenString, err := generateToken(user.ID, user.Email, s.jwtSecret)
	if err != nil {
		return nil, err
	}

	return &OAuthResp{
		LoginResp: LoginResp{
			Token: tokenString,
			User: struct {
				ID        string `json:"id"`
				Email     string `json:"email"`
				IsVisitor bool   `json:"isVisitor"`
			}{
				ID:        user.ID.Hex(),
				Email:     user.Email,
				IsVisitor: false,
			},
		},
	}, nil
}
