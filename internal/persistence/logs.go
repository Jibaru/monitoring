package persistence

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const logsCollectionName = "logs"

type Log struct {
	ID        primitive.ObjectID     `bson:"_id" json:"id"`
	AppID     primitive.ObjectID     `bson:"appId" json:"appId"`
	Timestamp time.Time              `bson:"timestamp" json:"timestamp"`
	Data      map[string]interface{} `bson:"data" json:"data"`
	Raw       string                 `bson:"raw" json:"raw"`
	Level     string                 `bson:"level" json:"level"`
}

func SaveLogs(ctx context.Context, db *mongo.Database, logs []Log) error {
	collection := db.Collection(logsCollectionName)
	_, err := collection.InsertMany(ctx, toAnySlice(logs), nil)
	return err
}

func ListLogs(ctx context.Context, db *mongo.Database, criteria Criteria) ([]Log, error) {
	collection := db.Collection(logsCollectionName)
	cursor, err := collection.Aggregate(ctx, criteria.MapToPipeline())
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	logs := make([]Log, 0)
	for cursor.Next(ctx) {
		var aLog Log
		if err := cursor.Decode(&aLog); err != nil {
			return nil, err
		}
		logs = append(logs, aLog)
	}

	return logs, nil
}
