package persistence

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"monitoring/internal/domain"
)

type logRepo struct {
	db         *mongo.Database
	collection string
}

type LogDoc struct {
	ID        primitive.ObjectID `bson:"_id"`
	AppID     primitive.ObjectID `bson:"appId"`
	Timestamp time.Time          `bson:"timestamp"`
	Data      map[string]any     `bson:"data"`
	Raw       string             `bson:"raw"`
	Level     string             `bson:"level"`
}

func logToDomain(log *LogDoc) (*domain.Log, error) {
	return domain.NewLog(
		log.ID,
		log.AppID,
		log.Timestamp,
		log.Data,
		log.Raw,
		log.Level,
	)
}

func logFromDomain(log domain.Log) LogDoc {
	return LogDoc{
		ID:        log.ID(),
		AppID:     log.AppID(),
		Timestamp: log.Timestamp(),
		Data:      log.Data(),
		Raw:       log.Raw(),
		Level:     log.Level(),
	}
}

func logsFromDomain(logs []domain.Log) []LogDoc {
	docs := make([]LogDoc, len(logs))
	for i, log := range logs {
		docs[i] = logFromDomain(log)
	}
	return docs
}

func NewLogRepo(db *mongo.Database) *logRepo {
	return &logRepo{db: db, collection: "logs"}
}

func (r *logRepo) SaveLogs(ctx context.Context, logs []domain.Log) error {
	collection := r.db.Collection(r.collection)
	_, err := collection.InsertMany(ctx, toAnySlice(logsFromDomain(logs)), nil)
	return err
}
func (r *logRepo) ListLogs(ctx context.Context, criteria domain.Criteria) ([]domain.Log, error) {
	collection := r.db.Collection(r.collection)
	cursor, err := collection.Aggregate(ctx, criteriaToPipeline(criteria))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	logs := make([]domain.Log, 0)
	for cursor.Next(ctx) {
		var aLog LogDoc
		if err := cursor.Decode(&aLog); err != nil {
			return nil, err
		}

		l, err := logToDomain(&aLog)
		if err != nil {
			return nil, err
		}

		logs = append(logs, *l)
	}

	return logs, nil
}

type Range struct {
	From time.Time
	To   time.Time
}

type LogSchemaResult struct {
	Total  int            `json:"total"`
	Schema map[string]int `json:"schema"`
}

func GetLogsSchema(ctx context.Context, db *mongo.Database, userID primitive.ObjectID, appIDs []primitive.ObjectID, optionalRange *Range) (LogSchemaResult, error) {
	logsColl := db.Collection("logs")

	pipeline := mongo.Pipeline{
		{
			{Key: "$lookup", Value: bson.M{
				"from":         "apps",
				"localField":   "appId",
				"foreignField": "_id",
				"as":           "app",
			}},
		},
		{{Key: "$unwind", Value: "$app"}},
	}

	matchFilter := bson.M{}
	if len(appIDs) > 0 {
		matchFilter["app._id"] = bson.M{"$in": appIDs}
	} else {
		matchFilter["app.userId"] = userID
	}
	if optionalRange != nil {
		matchFilter["timestamp"] = bson.M{
			"$gte": optionalRange.From,
			"$lte": optionalRange.To,
		}
	}
	pipeline = append(pipeline, bson.D{{Key: "$match", Value: matchFilter}})

	pipeline = append(pipeline, bson.D{{
		Key: "$facet", Value: bson.M{
			"total": []bson.M{
				{"$count": "total"},
			},
			"schema": []bson.M{
				{
					"$project": bson.M{
						"flattenedFields": bson.M{
							"$reduce": bson.M{
								"input": bson.M{
									"$map": bson.M{
										"input": bson.M{"$objectToArray": "$data"},
										"as":    "elem",
										"in": bson.M{
											"$cond": []interface{}{
												bson.M{"$eq": []interface{}{bson.M{"$type": "$$elem.v"}, "object"}},
												// If it is an object, an array is created with the subfields (level 2)
												bson.M{
													"$map": bson.M{
														"input": bson.M{"$objectToArray": "$$elem.v"},
														"as":    "sub",
														"in": bson.M{
															"$concat": []interface{}{"$$elem.k", ".", "$$sub.k"},
														},
													},
												},
												// Otherwise, an array with the key is returned.
												[]interface{}{"$$elem.k"},
											},
										},
									},
								},
								"initialValue": []interface{}{},
								"in":           bson.M{"$concatArrays": []interface{}{"$$value", "$$this"}},
							},
						},
					},
				},
				// Ensure that each document has a string in flattenedFields after unpacking.
				{"$unwind": "$flattenedFields"},
				{"$group": bson.M{
					"_id":   "$flattenedFields",
					"count": bson.M{"$sum": 1},
				}},
			},
		},
	}})

	cursor, err := logsColl.Aggregate(ctx, pipeline)
	if err != nil {
		return LogSchemaResult{}, err
	}
	defer cursor.Close(ctx)

	var aggregateResult []struct {
		Total []struct {
			Total int `bson:"total"`
		} `bson:"total"`
		Schema []struct {
			ID    string `bson:"_id"`
			Count int    `bson:"count"`
		} `bson:"schema"`
	}

	if err := cursor.All(ctx, &aggregateResult); err != nil {
		return LogSchemaResult{}, err
	}

	result := LogSchemaResult{
		Schema: make(map[string]int),
		Total:  0,
	}
	if len(aggregateResult) > 0 {
		if len(aggregateResult[0].Total) > 0 {
			result.Total = aggregateResult[0].Total[0].Total
		}
		for _, s := range aggregateResult[0].Schema {
			result.Schema[s.ID] = s.Count
		}
	}

	return result, nil
}
