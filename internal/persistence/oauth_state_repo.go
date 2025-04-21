package persistence

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"monitoring/internal/domain"
)

var _ domain.OAuthStateRepo = &oAuthStateRepo{}

type oAuthStateRepo struct {
	db         *mongo.Database
	collection string
}

type OAuthStateDoc struct {
	ID    primitive.ObjectID `bson:"_id"`
	State string             `bson:"state"`
}

func oAuthStateFromDomain(oAuthState domain.OAuthState) OAuthStateDoc {
	return OAuthStateDoc{
		ID:    oAuthState.ID(),
		State: oAuthState.State(),
	}
}

func NewOAuthStateRepo(db *mongo.Database) *oAuthStateRepo {
	return &oAuthStateRepo{
		db:         db,
		collection: "oauthState",
	}
}

func (r *oAuthStateRepo) SaveOAuthState(ctx context.Context, oauthState domain.OAuthState) error {
	collection := r.db.Collection(r.collection)
	_, err := collection.InsertOne(ctx, oAuthStateFromDomain(oauthState))
	return err
}

func (r *oAuthStateRepo) DeleteOAuthStateByState(ctx context.Context, state string) error {
	collection := r.db.Collection(r.collection)
	res, err := collection.DeleteOne(ctx, map[string]string{"state": state})
	if err != nil {
		return err
	}
	if res.DeletedCount == 0 {
		return domain.ErrNoOAuthStatesDeleted
	}
	return nil
}
