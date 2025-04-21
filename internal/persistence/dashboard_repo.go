package persistence

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"monitoring/internal/domain"
)

var _ domain.DashboardRepo = &dashboardRepo{}

type dashboardRepo struct {
	db *mongo.Database
}

func NewDashboardRepo(db *mongo.Database) *dashboardRepo {
	return &dashboardRepo{
		db: db,
	}
}

func (r *dashboardRepo) OverviewKPIs(
	ctx context.Context,
	userID domain.ID,
	optionalRange *domain.Range,
) (domain.DashboardOverviewKPIs, error) {
	logsColl := r.db.Collection("logs")

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
		return domain.DashboardOverviewKPIs{}, err
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
		return domain.DashboardOverviewKPIs{}, fmt.Errorf("Cursor.All: %w", err)
	}
	if len(raw) == 0 {
		return domain.DashboardOverviewKPIs{
			LogsPerApp: make(map[string]struct {
				AppName string
				Total   int64
			}),
			LogsByPeriod: make(map[domain.Period]struct {
				Total int64
			}),
		}, nil
	}
	rw := raw[0]

	totalLogs := int64(0)
	if len(rw.TotalCount) > 0 {
		totalLogs = rw.TotalCount[0].Total
	}

	counts := map[string]int64{}
	for _, lc := range rw.LevelCounts {
		counts[lc.Level] = lc.Count
	}

	overview := domain.DashboardOverviewKPIs{
		Logs: domain.DashboardOverviewKPI{Total: totalLogs, Percentage: func() float64 {
			if totalLogs == 0 {
				return 0
			}
			return 100
		}()},
		Errors: domain.DashboardOverviewKPI{Total: counts["ERROR"], Percentage: func() float64 {
			if totalLogs == 0 {
				return 0
			}
			return float64(counts["ERROR"]) * 100 / float64(totalLogs)
		}()},
		Warnings: domain.DashboardOverviewKPI{Total: counts["WARNING"], Percentage: func() float64 {
			if totalLogs == 0 {
				return 0
			}
			return float64(counts["WARNING"]) * 100 / float64(totalLogs)
		}()},
		Info: domain.DashboardOverviewKPI{Total: counts["INFO"], Percentage: func() float64 {
			if totalLogs == 0 {
				return 0
			}
			return float64(counts["INFO"]) * 100 / float64(totalLogs)
		}()},
		LogsPerApp: make(map[string]struct {
			AppName string
			Total   int64
		}),
		LogsByPeriod: make(map[domain.Period]struct {
			Total int64
		}),
	}

	for _, lp := range rw.LogsPerApp {
		overview.LogsPerApp[lp.ID] = struct {
			AppName string
			Total   int64
		}{AppName: lp.AppName, Total: lp.Total}
	}
	for _, p := range rw.LogsByPeriod {
		per := domain.Period{Year: p.ID.Year, Month: p.ID.Month}
		overview.LogsByPeriod[per] = struct {
			Total int64
		}{Total: p.Total}
	}

	return overview, nil
}
