package domain

import "context"

type UserRepo interface {
	SaveUser(ctx context.Context, user User) error
	ExistUserByEmail(ctx context.Context, email string) (bool, error)
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	GetUserByID(ctx context.Context, id ID) (*User, error)
	UpdateUser(ctx context.Context, user User) error
	ListUsers(ctx context.Context, criteria Criteria) ([]User, error)
	DeleteUser(ctx context.Context, id ID) error
}
