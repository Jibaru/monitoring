package persistence

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const appsCollectionName = "apps"

type App struct {
	ID        primitive.ObjectID `bson:"_id" json:"id"`
	Name      string             `bson:"name" json:"name"`
	AppKey    string             `bson:"appKey" json:"appKey"`
	UserID    string             `bson:"userId" json:"userId"`
	CreatedAt time.Time          `bson:"createdAt" json:"createdAt"`
}

func SaveApp(ctx context.Context, db *mongo.Database, app App) error {
	collection := db.Collection(appsCollectionName)
	_, err := collection.InsertOne(ctx, app)
	return err
}

func DeleteApp(ctx context.Context, db *mongo.Database, appID string) error {
	collection := db.Collection(appsCollectionName)
	_, err := collection.DeleteOne(ctx, map[string]string{"_id": appID})
	return err
}

func ListAppsPaginated(ctx context.Context, db *mongo.Database, userID string, page int, limit int) ([]App, error) {
	collection := db.Collection(appsCollectionName)
	cursor, err := collection.Aggregate(ctx, mongo.Pipeline{
		{{Key: "$match", Value: map[string]string{"userId": userID}}},
		{{Key: "$skip", Value: (page - 1) * limit}},
		{{Key: "$limit", Value: limit}},
	})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	apps := make([]App, 0)
	for cursor.Next(ctx) {
		var app App
		if err := cursor.Decode(&app); err != nil {
			return nil, err
		}
		apps = append(apps, app)
	}

	return apps, nil
}
