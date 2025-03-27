package scripts

import (
	"context"
	"fmt"
	"monitoring/internal/persistence"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ReceiveLogsReq struct {
	AppID  string                   `json:"-"`
	APIKey string                   `json:"-"`
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
	appID, err := primitive.ObjectIDFromHex(req.AppID)
	if err != nil {
		return nil, fmt.Errorf("ID de app inválido")
	}

	var app struct {
		AppKey string `bson:"app_key"`
	}
	if err := s.db.Collection("apps").FindOne(ctx, bson.M{"_id": req.AppID}).Decode(&app); err != nil {
		return nil, fmt.Errorf("app no encontrada")
	}
	if app.AppKey != req.APIKey {
		return nil, fmt.Errorf("API key inválida")
	}

	for _, logEntry := range req.Logs {
		logDoc := persistence.Log{
			ID:        primitive.NewObjectID(),
			AppID:     appID,
			Timestamp: time.Now().UTC(),
			Data:      logEntry,
		}

		err := persistence.SaveLog(ctx, s.db, logDoc)
		if err != nil {
			return nil, err
		}
	}
	return &ReceiveLogsResp{Message: "Logs recibidos"}, nil
}
