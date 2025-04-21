package scripts

import (
	"context"

	"monitoring/internal/domain"
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
	userRepo   domain.UserRepo
	mailSender *mail.MailSender
	webBaseURI string
}

func NewRegisterScript(
	userRepo domain.UserRepo,
	mailSender *mail.MailSender,
	webBaseURI string,
) *RegisterScript {
	return &RegisterScript{userRepo: userRepo, mailSender: mailSender, webBaseURI: webBaseURI}
}

func (s *RegisterScript) Exec(ctx context.Context, req RegisterReq) (*RegisterResp, error) {
	createUserScript := NewCreateUserScript(s.userRepo, s.mailSender, s.webBaseURI)
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
