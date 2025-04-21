package scripts

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"monitoring/internal/domain"
	"monitoring/internal/mail"
)

type CreateUserReq struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	VisitorID string `json:"visitorId"`
	RootID    string `json:"rootId"`
}

type CreateUserResp struct {
	ID    string `json:"id"`
	Email string `json:"email"`
}

type CreateUserScript struct {
	userRepo   domain.UserRepo
	mailSender *mail.MailSender
	webBaseURI string
}

func NewCreateUserScript(
	userRepo domain.UserRepo,
	mailSender *mail.MailSender,
	webBaseURI string,
) *CreateUserScript {
	return &CreateUserScript{userRepo: userRepo, mailSender: mailSender, webBaseURI: webBaseURI}
}

func (s *CreateUserScript) Exec(ctx context.Context, req CreateUserReq) (*CreateUserResp, error) {
	id := domain.NewAutoID()
	isFromVisitor := false
	if req.VisitorID != "" {
		visitorUserID, err := domain.NewID(req.VisitorID)
		if err != nil {
			return nil, err
		}

		visitorUser, err := s.userRepo.GetUserByID(ctx, visitorUserID)
		if err != nil {
			return nil, err
		}

		if !visitorUser.IsVisitor() {
			return nil, errors.New("visitor user is not visitor")
		}

		id = visitorUserID
		isFromVisitor = true
	}

	var rootUserID *domain.ID
	if req.RootID != "" {
		rootID, err := domain.NewID(req.RootID)
		if err != nil {
			return nil, err
		}
		rootUser, err := s.userRepo.GetUserByID(ctx, rootID)
		if err != nil {
			return nil, err
		}

		if !rootUser.IsRoot() {
			return nil, errors.New("root user is not root")
		}

		if rootUser.IsVisitor() {
			return nil, errors.New("root user is visitor")
		}

		rootUserID = &rootID
	}

	exists, err := s.userRepo.ExistUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}

	if exists {
		return nil, errors.New("user with this email already exists")
	}

	encryptedPassword, err := encryptPassword(req.Password)
	if err != nil {
		return nil, err
	}

	user, err := domain.NewUser(
		id,
		generateUsername(),
		req.Email,
		encryptedPassword,
		Now().UTC(),
		s.generatePin(),
		Now().UTC().Add(1*24*time.Hour),
		nil,
		false,
		false,
		rootUserID,
	)
	if err != nil {
		return nil, err
	}

	if isFromVisitor {
		err = s.userRepo.UpdateUser(ctx, *user)
		if err != nil {
			return nil, err
		}
	} else {
		err = s.userRepo.SaveUser(ctx, *user)
		if err != nil {
			return nil, err
		}
	}

	validatePinURL := s.webBaseURI + "/validate?userId=" + user.ID().Hex()

	err = s.mailSender.Send(
		req.Email,
		"Validate your account",
		fmt.Sprintf("Your pin is %v. You have to validate ir here: %s until %s", user.Pin(), validatePinURL, user.PinExpiresAt().Format(time.RFC822Z)),
	)
	if err != nil {
		return nil, err
	}

	return &CreateUserResp{ID: user.ID().Hex(), Email: req.Email}, nil
}

func (s *CreateUserScript) generatePin() string {
	const length = 6
	letters := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	digits := "0123456789"
	var pin []byte

	// At least 2 letters
	for i := 0; i < 2; i++ {
		pin = append(pin, letters[rand.Intn(len(letters))])
	}
	// At least 2 numbers
	for i := 0; i < 2; i++ {
		pin = append(pin, digits[rand.Intn(len(digits))])
	}

	allChars := letters + digits
	for i := 0; i < length-4; i++ {
		pin = append(pin, allChars[rand.Intn(len(allChars))])
	}

	rand.Shuffle(len(pin), func(i, j int) {
		pin[i], pin[j] = pin[j], pin[i]
	})

	return string(pin)
}
