package persistence

import (
	"fmt"
	"monitoring/internal/domain"

	"go.mongodb.org/mongo-driver/bson"
)

func toAnySlice[T any](values []T) []any {
	docs := make([]any, len(values))
	for i, value := range values {
		docs[i] = value
	}
	return docs
}

// criteriaToPipeline maps the criteria to a MongoDB aggregation pipeline.
func criteriaToPipeline(c domain.Criteria) []bson.M {
	var pipeline []bson.M

	if len(c.Filters) > 0 {
		filterStage := bson.M{}
		for _, f := range c.Filters {
			switch f.Type {
			case domain.Equals:
				filterStage[f.Field] = bson.M{"$eq": f.Value}
			case domain.NotEquals:
				filterStage[f.Field] = bson.M{"$ne": f.Value}
			case domain.Like:
				pattern := fmt.Sprintf(".*%v.*", f.Value)
				filterStage[f.Field] = bson.M{"$regex": pattern, "$options": "i"}
			case domain.In:
				filterStage[f.Field] = bson.M{"$in": f.Value}
			case domain.GreaterThanOrEqual:
				filterStage[f.Field] = bson.M{"$gte": f.Value}
			case domain.LessThanOrEqual:
				filterStage[f.Field] = bson.M{"$lte": f.Value}
			}
		}
		pipeline = append(pipeline, bson.M{"$match": filterStage})
	}

	if c.Sort.Field != "" {
		order := 1
		if c.Sort.Order == domain.Desc {
			order = -1
		}
		pipeline = append(pipeline, bson.M{"$sort": bson.M{c.Sort.Field: order}})
	}

	if c.Pagination.Offset > 0 {
		pipeline = append(pipeline, bson.M{"$skip": c.Pagination.Offset})
	}
	if c.Pagination.Limit > 0 {
		pipeline = append(pipeline, bson.M{"$limit": c.Pagination.Limit})
	}

	return pipeline
}
