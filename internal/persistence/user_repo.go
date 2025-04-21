package persistence

import (
	"context"
	"errors"
	"monitoring/internal/domain"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserDoc struct {
	ID           primitive.ObjectID  `bson:"_id" json:"id"`
	Username     string              `bson:"username" json:"username"`
	Email        string              `bson:"email" json:"email"`
	Password     string              `bson:"password" json:"password"`
	RegisteredAt time.Time           `bson:"registeredAt" json:"registeredAt"`
	Pin          string              `bson:"pin" json:"pin"`
	PinExpiresAt time.Time           `bson:"pinExpiresAt" json:"pinExpiresAt"`
	ValidatedAt  *time.Time          `bson:"validatedAt" json:"validatedAt"`
	IsVisitor    bool                `bson:"isVisitor" json:"isVisitor"`
	FromOAuth    bool                `bson:"fromOAuth" json:"fromOAuth"`
	RootUserID   *primitive.ObjectID `bson:"rootUserId,omitempty" json:"rootUserId"`
}

type userRepo struct {
	db         *mongo.Database
	collection string
}

func userFromDomain(user domain.User) UserDoc {
	return UserDoc{
		ID:           user.ID(),
		Username:     user.Username(),
		Email:        user.Email(),
		Password:     user.Password(),
		RegisteredAt: user.RegisteredAt(),
		Pin:          user.Pin(),
		PinExpiresAt: user.PinExpiresAt(),
		ValidatedAt:  user.ValidatedAt(),
		IsVisitor:    user.IsVisitor(),
		FromOAuth:    user.FromOAuth(),
		RootUserID:   user.RootUserID(),
	}
}

func userToDomain(user *UserDoc) (*domain.User, error) {
	return domain.NewUser(
		user.ID,
		user.Username,
		user.Email,
		user.Password,
		user.RegisteredAt,
		user.Pin,
		user.PinExpiresAt,
		user.ValidatedAt,
		user.IsVisitor,
		user.FromOAuth,
		user.RootUserID,
	)
}

func NewUserRepo(db *mongo.Database) *userRepo {
	return &userRepo{
		db:         db,
		collection: "users",
	}
}

func (r *userRepo) SaveUser(ctx context.Context, user domain.User) error {
	collection := r.db.Collection(r.collection)
	_, err := collection.InsertOne(ctx, userFromDomain(user))
	return err
}

func (r *userRepo) ExistUserByEmail(ctx context.Context, email string) (bool, error) {
	var user UserDoc
	err := r.db.Collection(r.collection).FindOne(ctx, map[string]string{"email": email}).Decode(&user)
	if err != nil && errors.Is(err, mongo.ErrNoDocuments) {
		return false, nil
	} else if err != nil {
		return false, err
	}
	return true, nil
}

func (r *userRepo) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	var user UserDoc
	err := r.db.Collection(r.collection).FindOne(ctx, map[string]string{"email": email}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return userToDomain(&user)
}

func (r *userRepo) GetUserByID(ctx context.Context, id domain.ID) (*domain.User, error) {
	var user UserDoc
	err := r.db.Collection(r.collection).FindOne(ctx, map[string]any{"_id": id}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return userToDomain(&user)
}

func (r *userRepo) UpdateUser(ctx context.Context, user domain.User) error {
	collection := r.db.Collection(r.collection)
	_, err := collection.UpdateOne(ctx, bson.M{"_id": user.ID()}, map[string]any{
		"$set": userFromDomain(user),
	})
	return err
}
