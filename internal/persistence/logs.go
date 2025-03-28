package persistence

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const logsCollectionName = "logs"

type Log struct {
	ID        primitive.ObjectID     `bson:"_id" json:"id"`
	AppID     primitive.ObjectID     `bson:"appId" json:"appId"`
	Timestamp time.Time              `bson:"timestamp" json:"timestamp"`
	Data      map[string]interface{} `bson:"data" json:"data"`
}

func SaveLog(ctx context.Context, db *mongo.Database, log Log) error {
	collection := db.Collection(logsCollectionName)
	_, err := collection.InsertOne(ctx, log)
	return err
}

func SearchLogs(ctx context.Context, db *mongo.Database, appID string, from time.Time, to time.Time, page int, limit int, extraMatch bson.M) ([]Log, error) {
	collection := db.Collection(logsCollectionName)
	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.M{
			"appId": appID,
			"timestamp": bson.M{
				"$gte": from,
				"$lte": to,
			},
		}}},
		{{Key: "$sort", Value: bson.M{"timestamp": -1}}},
		{{Key: "$skip", Value: (page - 1) * limit}},
		{{Key: "$limit", Value: limit}},
	}
	if extraMatch != nil {
		pipeline = append(pipeline, bson.D{{Key: "$match", Value: extraMatch}})
	}

	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	var logs []Log
	for cursor.Next(ctx) {
		var log Log
		if err := cursor.Decode(&log); err != nil {
			return nil, err
		}
		logs = append(logs, log)
	}

	return logs, nil
}
