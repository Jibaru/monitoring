package persistence

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"monitoring/internal/domain"
)

var _ domain.LogSchemaRepo = &logSchemaRepo{}

type logSchemaRepo struct {
	db *mongo.Database
}

func NewLogSchemaRepo(db *mongo.Database) *logSchemaRepo {
	return &logSchemaRepo{
		db: db,
	}
}

func (r *logSchemaRepo) Get(ctx context.Context, userID domain.ID, appIDs []domain.ID, optionalRange *domain.Range) (domain.LogSchema, error) {
	logsColl := r.db.Collection("logs")

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
		return domain.LogSchema{}, err
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
		return domain.LogSchema{}, err
	}

	result := domain.LogSchema{
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
