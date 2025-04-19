package scripts

import (
	"context"
	"fmt"
	"monitoring/internal/persistence"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type VisitorLoginResp struct {
	LoginResp
}

type VisitorLoginScript struct {
	db        *mongo.Database
	jwtSecret []byte
}

func NewVisitorLoginScript(db *mongo.Database, jwtSecret string) *VisitorLoginScript {
	return &VisitorLoginScript{
		db:        db,
		jwtSecret: []byte(jwtSecret),
	}
}

func (s *VisitorLoginScript) Exec(ctx context.Context) (*VisitorLoginResp, error) {
	visitorEmail := uuid.NewString() + "_monitoring_" + fmt.Sprintf("%v", time.Now().Unix()) + "@mail.app"

	user := persistence.User{
		ID:           primitive.NewObjectID(),
		Username:     generateUsername(),
		Email:        visitorEmail,
		Password:     "",
		RegisteredAt: time.Now().UTC(),
		ValidatedAt:  nil,
		IsVisitor:    true,
		RootUserID:   nil,
	}

	err := persistence.SaveUser(ctx, s.db, user)
	if err != nil {
		return nil, err
	}

	tokenString, err := generateToken(user.ID, user.Email, s.jwtSecret)
	if err != nil {
		return nil, err
	}

	return &VisitorLoginResp{
		LoginResp: *userToLoginResp(tokenString, &user),
	}, nil
}
