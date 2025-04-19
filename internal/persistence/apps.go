package persistence

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"monitoring/internal/domain"
)

var _ domain.AppRepo = &appRepo{}

type appRepo struct {
	db         *mongo.Database
	collection string
}

type AppDoc struct {
	ID        primitive.ObjectID `bson:"_id"`
	Name      string             `bson:"name"`
	AppKey    string             `bson:"appKey"`
	UserID    primitive.ObjectID `bson:"userId"`
	CreatedAt time.Time          `bson:"createdAt"`
}

func appFromDomain(app domain.App) AppDoc {
	return AppDoc{
		ID:        app.ID(),
		Name:      app.Name(),
		AppKey:    app.AppKey(),
		UserID:    app.UserID(),
		CreatedAt: app.CreatedAt(),
	}
}

func appToDomain(app *AppDoc) (*domain.App, error) {
	return domain.NewApp(
		app.ID,
		app.Name,
		app.AppKey,
		app.UserID,
		app.CreatedAt,
	)
}

func NewAppRepo(db *mongo.Database) *appRepo {
	return &appRepo{db: db, collection: "apps"}
}

func (r *appRepo) SaveApp(ctx context.Context, app domain.App) error {
	collection := r.db.Collection(r.collection)
	_, err := collection.InsertOne(ctx, appFromDomain(app))
	return err
}

func (r *appRepo) UpdateApp(ctx context.Context, app domain.App) error {
	collection := r.db.Collection(r.collection)
	_, err := collection.UpdateOne(ctx, bson.M{"_id": app.ID}, map[string]any{
		"$set": appFromDomain(app),
	})
	return err
}

func (r *appRepo) GetAppByID(ctx context.Context, appID domain.ID) (*domain.App, error) {
	var app AppDoc
	err := r.db.Collection(r.collection).FindOne(ctx, bson.M{"_id": appID}).Decode(&app)
	if err != nil {
		return nil, err
	}
	return appToDomain(&app)
}

func (r *appRepo) GetAppByKey(ctx context.Context, appKey string) (*domain.App, error) {
	var app AppDoc
	err := r.db.Collection(r.collection).FindOne(ctx, bson.M{"appKey": appKey}).Decode(&app)
	if err != nil {
		return nil, err
	}
	return appToDomain(&app)
}

func (r *appRepo) DeleteApp(ctx context.Context, appID domain.ID) error {
	collection := r.db.Collection(r.collection)
	_, err := collection.DeleteOne(ctx, map[string]any{"_id": appID})
	return err
}

func (r *appRepo) ListApps(ctx context.Context, criteria domain.Criteria) ([]domain.App, error) {
	collection := r.db.Collection(r.collection)
	cursor, err := collection.Aggregate(ctx, criteriaToPipeline(criteria))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	apps := make([]domain.App, 0)
	for cursor.Next(ctx) {
		var app AppDoc
		if err := cursor.Decode(&app); err != nil {
			return nil, err
		}

		domainApp, err := appToDomain(&app)
		if err != nil {
			return nil, err
		}

		apps = append(apps, *domainApp)
	}

	return apps, nil
}
