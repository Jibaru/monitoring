package scripts

import (
	"context"
	"errors"
	"monitoring/internal/persistence"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type RegisterReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterResp struct {
	ID    string `json:"id"`
	Email string `json:"email"`
}

type RegisterScript struct {
	db *mongo.Database
}

func NewRegisterScript(db *mongo.Database) *RegisterScript {
	return &RegisterScript{db: db}
}

func (s *RegisterScript) Exec(ctx context.Context, req RegisterReq) (*RegisterResp, error) {
	exists, err := persistence.ExistUserByEmail(ctx, s.db, req.Email)
	if err != nil {
		return nil, err
	}

	if exists {
		return nil, errors.New("el usuario con este email ya existe")
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	user := persistence.User{
		ID:           primitive.NewObjectID(),
		Email:        req.Email,
		Password:     string(hashed),
		RegisteredAt: time.Now().UTC(),
	}

	err = persistence.SaveUser(ctx, s.db, user)
	if err != nil {
		return nil, err
	}

	return &RegisterResp{ID: user.ID.Hex(), Email: req.Email}, nil
}
