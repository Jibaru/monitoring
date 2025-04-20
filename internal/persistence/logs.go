package persistence

import (
	"context"
	"fmt"
	"monitoring/internal/domain"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type logRepo struct {
	db         *mongo.Database
	collection string
}

type LogDoc struct {
	ID        primitive.ObjectID     `bson:"_id" json:"id"`
	AppID     primitive.ObjectID     `bson:"appId" json:"appId"`
	Timestamp time.Time              `bson:"timestamp" json:"timestamp"`
	Data      map[string]interface{} `bson:"data" json:"data"`
	Raw       string                 `bson:"raw" json:"raw"`
	Level     string                 `bson:"level" json:"level"`
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

type DashboardOverviewKPI struct {
	Total      int64   `bson:"total"`
	Percentage float64 `bson:"percentage"`
}

type Period struct {
	Year  int `bson:"year"`
	Month int `bson:"month"`
}

type DashboardOverviewKPIs struct {
	Logs       DashboardOverviewKPI `bson:"logs"`
	Errors     DashboardOverviewKPI `bson:"errors"`
	Warnings   DashboardOverviewKPI `bson:"warnings"`
	Info       DashboardOverviewKPI `bson:"info"`
	LogsPerApp map[string]struct {
		AppName string `bson:"appName"`
		Total   int64  `bson:"total"`
	} `bson:"logsPerApp"`
	LogsByPeriod map[Period]struct {
		Total int64 `bson:"total"`
	} `bson:"logsByPeriod"`
}

type Range struct {
	From time.Time
	To   time.Time
}

func GetDashboardOverviewKPIs(
	ctx context.Context,
	db *mongo.Database,
	userID primitive.ObjectID,
	optionalRange *Range,
) (DashboardOverviewKPIs, error) {
	logsColl := db.Collection("logs")

	pipeline := mongo.Pipeline{}

	if optionalRange != nil {
		pipeline = append(pipeline, bson.D{{Key: "$match", Value: bson.D{
			{Key: "timestamp", Value: bson.D{
				{Key: "$gte", Value: optionalRange.From},
				{Key: "$lte", Value: optionalRange.To},
			}},
		}}})
	}

	pipeline = append(pipeline,
		bson.D{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "apps"},
			{Key: "localField", Value: "appId"},
			{Key: "foreignField", Value: "_id"},
			{Key: "as", Value: "app"},
		}}},
		bson.D{{Key: "$unwind", Value: "$app"}},
		bson.D{{Key: "$match", Value: bson.D{
			{Key: "app.userId", Value: userID},
		}}},
	)

	pipeline = append(pipeline, bson.D{{Key: "$facet", Value: bson.D{
		{Key: "logsPerApp", Value: bson.A{
			bson.D{{Key: "$group", Value: bson.D{
				{Key: "_id", Value: "$app._id"},
				{Key: "appName", Value: bson.D{{Key: "$first", Value: "$app.name"}}},
				{Key: "total", Value: bson.D{{Key: "$sum", Value: 1}}},
			}}},
		}},
		{Key: "totalCount", Value: bson.A{
			bson.D{{Key: "$count", Value: "total"}},
		}},
		{Key: "levelCounts", Value: bson.A{
			bson.D{{Key: "$group", Value: bson.D{
				{Key: "_id", Value: "$level"},
				{Key: "count", Value: bson.D{{Key: "$sum", Value: 1}}},
			}}},
		}},
		{Key: "logsByPeriod", Value: bson.A{
			bson.D{{Key: "$group", Value: bson.D{
				{Key: "_id", Value: bson.D{
					{Key: "year", Value: bson.D{{Key: "$year", Value: "$timestamp"}}},
					{Key: "month", Value: bson.D{{Key: "$month", Value: "$timestamp"}}},
				}},
				{Key: "total", Value: bson.D{{Key: "$sum", Value: 1}}},
			}}},
		}},
	}}})

	cur, err := logsColl.Aggregate(ctx, pipeline)
	if err != nil {
		return DashboardOverviewKPIs{}, err
	}
	defer cur.Close(ctx)

	var raw []struct {
		LogsPerApp []struct {
			ID      string `bson:"_id"`
			AppName string `bson:"appName"`
			Total   int64  `bson:"total"`
		} `bson:"logsPerApp"`
		TotalCount []struct {
			Total int64 `bson:"total"`
		} `bson:"totalCount"`
		LevelCounts []struct {
			Level string `bson:"_id"`
			Count int64  `bson:"count"`
		} `bson:"levelCounts"`
		LogsByPeriod []struct {
			ID struct {
				Year  int `bson:"year"`
				Month int `bson:"month"`
			} `bson:"_id"`
			Total int64 `bson:"total"`
		} `bson:"logsByPeriod"`
	}
	if err := cur.All(ctx, &raw); err != nil {
		return DashboardOverviewKPIs{}, fmt.Errorf("Cursor.All: %w", err)
	}
	if len(raw) == 0 {
		return DashboardOverviewKPIs{
			LogsPerApp: make(map[string]struct {
				AppName string "bson:\"appName\""
				Total   int64  "bson:\"total\""
			}),
			LogsByPeriod: make(map[Period]struct {
				Total int64 "bson:\"total\""
			}),
		}, nil
	}
	r := raw[0]

	// 6) Total de logs
	totalLogs := int64(0)
	if len(r.TotalCount) > 0 {
		totalLogs = r.TotalCount[0].Total
	}

	// 7) Mapeo conteos por nivel
	counts := map[string]int64{}
	for _, lc := range r.LevelCounts {
		counts[lc.Level] = lc.Count
	}

	// 8) Ensamblo el resultado final
	overview := DashboardOverviewKPIs{
		Logs: DashboardOverviewKPI{Total: totalLogs, Percentage: func() float64 {
			if totalLogs == 0 {
				return 0
			}
			return 100
		}()},
		Errors: DashboardOverviewKPI{Total: counts["ERROR"], Percentage: func() float64 {
			if totalLogs == 0 {
				return 0
			}
			return float64(counts["ERROR"]) * 100 / float64(totalLogs)
		}()},
		Warnings: DashboardOverviewKPI{Total: counts["WARNING"], Percentage: func() float64 {
			if totalLogs == 0 {
				return 0
			}
			return float64(counts["WARNING"]) * 100 / float64(totalLogs)
		}()},
		Info: DashboardOverviewKPI{Total: counts["INFO"], Percentage: func() float64 {
			if totalLogs == 0 {
				return 0
			}
			return float64(counts["INFO"]) * 100 / float64(totalLogs)
		}()},
		LogsPerApp: make(map[string]struct {
			AppName string "bson:\"appName\""
			Total   int64  "bson:\"total\""
		}),
		LogsByPeriod: make(map[Period]struct {
			Total int64 "bson:\"total\""
		}),
	}

	for _, lp := range r.LogsPerApp {
		overview.LogsPerApp[lp.ID] = struct {
			AppName string `bson:"appName"`
			Total   int64  `bson:"total"`
		}{AppName: lp.AppName, Total: lp.Total}
	}
	for _, p := range r.LogsByPeriod {
		per := Period{Year: p.ID.Year, Month: p.ID.Month}
		overview.LogsByPeriod[per] = struct {
			Total int64 `bson:"total"`
		}{Total: p.Total}
	}

	return overview, nil
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
