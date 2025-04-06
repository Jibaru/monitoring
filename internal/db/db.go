package db

import (
	"context"
	"log"
	"monitoring/config"
	"time"

	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func New(cfg config.Config) (*mongo.Database, *mongo.Client) {
	cmdMonitor := &event.CommandMonitor{
		Started: func(_ context.Context, evt *event.CommandStartedEvent) {
			log.Print(evt.Command)
		},
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.MongoURI).SetMonitor(cmdMonitor))
	if err != nil {
		log.Fatal(err)
	}

	return client.Database(cfg.DBName), client
}
