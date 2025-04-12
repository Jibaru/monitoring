package persistence

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const usersCollectionName = "users"

type User struct {
	ID           primitive.ObjectID `bson:"_id" json:"id"`
	Email        string             `bson:"email" json:"email"`
	Password     string             `bson:"password" json:"password"`
	RegisteredAt time.Time          `bson:"registeredAt" json:"registeredAt"`
	Pin          string             `bson:"pin" json:"pin"`
	PinExpiresAt time.Time          `bson:"pinExpiresAt" json:"pinExpiresAt"`
	ValidatedAt  *time.Time         `bson:"validatedAt" json:"validatedAt"`
	IsVisitor    bool               `bson:"isVisitor" json:"isVisitor"`
	FromGithub   bool               `bson:"fromGithub" json:"fromGithub"`
}

func SaveUser(ctx context.Context, db *mongo.Database, user User) error {
	collection := db.Collection(usersCollectionName)
	_, err := collection.InsertOne(ctx, user)
	return err
}

func ExistUserByEmail(ctx context.Context, db *mongo.Database, email string) (bool, error) {
	var user User
	err := db.Collection(usersCollectionName).FindOne(ctx, map[string]string{"email": email}).Decode(&user)
	if err != nil && errors.Is(err, mongo.ErrNoDocuments) {
		return false, nil
	} else if err != nil {
		return false, err
	}
	return true, nil
}

func GetUserByEmail(ctx context.Context, db *mongo.Database, email string) (*User, error) {
	var user User
	err := db.Collection(usersCollectionName).FindOne(ctx, map[string]string{"email": email}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func GetUserByID(ctx context.Context, db *mongo.Database, id primitive.ObjectID) (*User, error) {
	var user User
	err := db.Collection(usersCollectionName).FindOne(ctx, map[string]any{"_id": id}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func UpdateUser(ctx context.Context, db *mongo.Database, user User) error {
	collection := db.Collection(usersCollectionName)
	_, err := collection.UpdateOne(ctx, bson.M{"_id": user.ID}, map[string]any{
		"$set": user,
	})
	return err
}
