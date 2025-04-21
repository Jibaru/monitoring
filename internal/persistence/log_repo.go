package persistence

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"monitoring/internal/domain"
)

var _ domain.LogRepo = &logRepo{}

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
