package scripts

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"monitoring/internal/persistence"
)

type ReceiveLogsReq struct {
	AppKey string                   `json:"-"`
	Logs   []map[string]interface{} `json:"logs"`
}

type ReceiveLogsResp struct {
	Message string `json:"message"`
}

type ReceiveLogsScript struct {
	db *mongo.Database
}

func NewReceiveLogsScript(db *mongo.Database) *ReceiveLogsScript {
	return &ReceiveLogsScript{db: db}
}

func (s *ReceiveLogsScript) Exec(ctx context.Context, req ReceiveLogsReq) (*ReceiveLogsResp, error) {
	app, err := persistence.GetAppByKey(ctx, s.db, req.AppKey)
	if err != nil {
		return nil, err
	}

	logs := make([]persistence.Log, len(req.Logs))
	for i, logEntry := range req.Logs {
		logs[i] = persistence.Log{
			ID:        primitive.NewObjectID(),
			AppID:     app.ID,
			Timestamp: time.Now().UTC(),
			Data:      logEntry,
		}
	}

	err = persistence.SaveLogs(ctx, s.db, logs)
	if err != nil {
		return nil, err
	}

	return &ReceiveLogsResp{Message: "Logs recibidos"}, nil
}
