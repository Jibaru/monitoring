package persistence

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const oauthStateCollectionName = "oauthState"

type OAuthState struct {
	ID    primitive.ObjectID `bson:"_id" json:"id"`
	State string             `bson:"state" json:"state"`
}

func SaveOAuthState(ctx context.Context, db *mongo.Database, oauthState OAuthState) error {
	collection := db.Collection(oauthStateCollectionName)
	_, err := collection.InsertOne(ctx, oauthState)
	return err
}

func ExistOAuthStateByState(ctx context.Context, db *mongo.Database, state string) (bool, error) {
	var oauthState OAuthState
	err := db.Collection(oauthStateCollectionName).FindOne(ctx, map[string]string{"sttate": state}).Decode(&oauthState)
	if err != nil && errors.Is(err, mongo.ErrNoDocuments) {
		return false, nil
	} else if err != nil {
		return false, err
	}
	return true, nil
}
