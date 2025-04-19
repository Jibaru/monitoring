package domain

import (
	"time"
)

type User struct {
	id           ID
	username     string
	email        string
	password     string
	registeredAt time.Time
	pin          string
	pinExpiresAt time.Time
	validatedAt  *time.Time
	isVisitor    bool
	fromOAuth    bool
	rootUserID   *ID
}

func NewUser(
	id ID,
	username string,
	email string,
	password string,
	registeredAt time.Time,
	pin string,
	pinExpiresAt time.Time,
	validatedAt *time.Time,
	isVisitor bool,
	fromOAuth bool,
	rootUserID *ID,
) (*User, error) {
	return &User{
		id:           id,
		username:     username,
		email:        email,
		password:     password,
		registeredAt: registeredAt,
		pin:          pin,
		pinExpiresAt: pinExpiresAt,
		validatedAt:  validatedAt,
		isVisitor:    isVisitor,
		fromOAuth:    fromOAuth,
		rootUserID:   rootUserID,
	}, nil
}
