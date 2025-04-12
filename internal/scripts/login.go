package scripts

import (
	"context"
	"errors"
	"monitoring/internal/persistence"

	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
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
	} `json:"user"`
}

type LoginScript struct {
	db        *mongo.Database
	jwtSecret []byte
}

func NewLoginScript(db *mongo.Database, jwtSecret string) *LoginScript {
	return &LoginScript{
		db:        db,
		jwtSecret: []byte(jwtSecret),
	}
}

func (s *LoginScript) Exec(ctx context.Context, req LoginReq) (*LoginResp, error) {
	user, err := persistence.GetUserByEmail(ctx, s.db, req.Email)
	if err != nil {
		return nil, err
	}

	if user.ValidatedAt == nil {
		return nil, errors.New("user should be validated to login")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, errors.New("email or password are invalid")
	}

	tokenString, err := generateToken(user.ID, user.Email, s.jwtSecret)
	if err != nil {
		return nil, err
	}

	return userToLoginResp(tokenString, user), nil
}

func userToLoginResp(token string, user *persistence.User) *LoginResp {
	return &LoginResp{
		Token: token,
		User: struct {
			ID        string `json:"id"`
			Username  string `json:"username"`
			Email     string `json:"email"`
			IsVisitor bool   `json:"isVisitor"`
		}{
			ID:        user.ID.Hex(),
			Username:  user.Username,
			Email:     user.Email,
			IsVisitor: user.IsVisitor,
		},
	}
}
