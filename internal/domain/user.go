package domain

import (
	"encoding/json"
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

func (u *User) ID() ID {
	return u.id
}

func (u *User) Username() string {
	return u.username
}

func (u *User) Email() string {
	return u.email
}

func (u *User) Password() string {
	return u.password
}

func (u *User) RegisteredAt() time.Time {
	return u.registeredAt
}

func (u *User) Pin() string {
	return u.pin
}

func (u *User) PinExpiresAt() time.Time {
	return u.pinExpiresAt
}

func (u *User) ValidatedAt() *time.Time {
	return u.validatedAt
}

func (u *User) IsVisitor() bool {
	return u.isVisitor
}

func (u *User) FromOAuth() bool {
	return u.fromOAuth
}

func (u *User) RootUserID() *ID {
	return u.rootUserID
}

func (u *User) IsRoot() bool {
	return u.rootUserID == nil
}

func (u *User) ChangeUsername(username string) error {
	u.username = username
	return nil
}

func (u *User) ChangeValidatedAt(validatedAt *time.Time) error {
	u.validatedAt = validatedAt
	return nil
}

func (u *User) ChangePassword(password string) error {
	u.password = password
	return nil
}

func (u User) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]any{
		"id":           u.id,
		"username":     u.username,
		"email":        u.email,
		"registeredAt": u.registeredAt,
		"pin":          u.pin,
		"pinExpiresAt": u.pinExpiresAt,
		"validatedAt":  u.validatedAt,
		"isVisitor":    u.isVisitor,
		"fromOAuth":    u.fromOAuth,
		"rootUserId":   u.rootUserID,
	})
}
