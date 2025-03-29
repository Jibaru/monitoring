package persistence

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
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

func UpdateApp(ctx context.Context, db *mongo.Database, app App) error {
	collection := db.Collection(appsCollectionName)
	_, err := collection.UpdateOne(ctx, bson.M{"_id": app.ID}, map[string]any{
		"$set": app,
	})
	return err
}

func GetAppByID(ctx context.Context, db *mongo.Database, appID primitive.ObjectID) (*App, error) {
	var app App
	err := db.Collection(appsCollectionName).FindOne(ctx, bson.M{"_id": appID}).Decode(&app)
	if err != nil {
		return nil, err
	}
	return &app, nil
}

func GetAppByKey(ctx context.Context, db *mongo.Database, appKey string) (*App, error) {
	var app App
	err := db.Collection(appsCollectionName).FindOne(ctx, bson.M{"appKey": appKey}).Decode(&app)
	if err != nil {
		return nil, err
	}
	return &app, nil
}

func DeleteApp(ctx context.Context, db *mongo.Database, appID string) error {
	collection := db.Collection(appsCollectionName)
	_, err := collection.DeleteOne(ctx, map[string]string{"_id": appID})
	return err
}

func ListApps(ctx context.Context, db *mongo.Database, criteria Criteria) ([]App, error) {
	collection := db.Collection(appsCollectionName)
	cursor, err := collection.Aggregate(ctx, criteria.MapToPipeline())
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
