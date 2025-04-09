package scripts

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"

	"monitoring/internal/mail"
	"monitoring/internal/persistence"
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
		ValidatedAt:  nil,
		Pin:          s.generatePin(),
		PinExpiresAt: time.Now().UTC().Add(1 * 24 * time.Hour),
	}

	err = persistence.SaveUser(ctx, s.db, user)
	if err != nil {
		return nil, err
	}

	validatePinURL := s.webBaseURI + "/validate?userId=" + user.ID.Hex()

	err = s.mailSender.Send(
		req.Email,
		"Validate your account",
		fmt.Sprintf("Your pin is %v. You have to validate ir here: %s until %s", user.Pin, validatePinURL, user.PinExpiresAt.Format(time.RFC822Z)),
	)
	if err != nil {
		return nil, err
	}

	return &RegisterResp{ID: user.ID.Hex(), Email: req.Email}, nil
}

func (s *RegisterScript) generatePin() string {
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
