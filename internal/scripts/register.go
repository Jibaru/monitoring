package scripts

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"

	"monitoring/internal/mail"
)

type RegisterReq struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	VisitorID string `json:"visitorId"`
}

type RegisterResp struct {
	ID    string `json:"id"`
	Email string `json:"email"`
}

type RegisterScript struct {
	db         *mongo.Database
	mailSender *mail.MailSender
	webBaseURI string
}

func NewRegisterScript(
	db *mongo.Database,
	mailSender *mail.MailSender,
	webBaseURI string,
) *RegisterScript {
	return &RegisterScript{db: db, mailSender: mailSender, webBaseURI: webBaseURI}
}

func (s *RegisterScript) Exec(ctx context.Context, req RegisterReq) (*RegisterResp, error) {
	createUserScript := NewCreateUserScript(s.db, s.mailSender, s.webBaseURI)
	resp, err := createUserScript.Exec(ctx, CreateUserReq{
		Email:     req.Email,
		Password:  req.Password,
		VisitorID: req.VisitorID,
		RootID:    "",
	})
	if err != nil {
		return nil, err
	}

	return &RegisterResp{ID: resp.ID, Email: resp.Email}, nil
}
