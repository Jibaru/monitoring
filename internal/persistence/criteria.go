package persistence

import "go.mongodb.org/mongo-driver/bson"

// FilterType represents the type of filter (equals or not equals).
type FilterType string

const (
	Equals    FilterType = "eq"
	NotEquals FilterType = "ne"
)

// Filter represents a filter on a specific field.
type Filter struct {
	Field string
	Type  FilterType
	Value any
}

// NewFilter creates a new filter.
func NewFilter(field string, filterType FilterType, value any) Filter {
	return Filter{Field: field, Type: filterType, Value: value}
}

// Pagination represents pagination settings with limit and offset.
type Pagination struct {
	Limit  int
	Offset int
}

// NewPagination creates a new pagination configuration.
func NewPagination(limit, offset int) Pagination {
	return Pagination{Limit: limit, Offset: offset}
}

// SortOrder represents the sorting order (ascending or descending).
type SortOrder string

const (
	Asc  SortOrder = "asc"
	Desc SortOrder = "desc"
)

// Sort represents sorting by a specific field.
type Sort struct {
	Field string
	Order SortOrder
}

// NewSort creates a new sorting configuration.
func NewSort(field string, order SortOrder) Sort {
	return Sort{Field: field, Order: order}
}

// Criteria represents all search criteria.
type Criteria struct {
	Filters    []Filter
	Pagination Pagination
	Sort       Sort
}

// NewCriteria creates a new set of search criteria.
func NewCriteria(filters []Filter, pagination Pagination, sort Sort) Criteria {
	return Criteria{Filters: filters, Pagination: pagination, Sort: sort}
}

// MapToPipeline maps the criteria to a MongoDB aggregation pipeline.
func (c Criteria) MapToPipeline() []bson.M {
	var pipeline []bson.M

	// Build filters
	if len(c.Filters) > 0 {
		filterStage := bson.M{}
		for _, f := range c.Filters {
			mongoOp := "$eq"
			if f.Type == NotEquals {
				mongoOp = "$ne"
			}
			filterStage[f.Field] = bson.M{mongoOp: f.Value}
		}
		pipeline = append(pipeline, bson.M{"$match": filterStage})
	}

	// Apply sorting
	if c.Sort.Field != "" {
		order := 1
		if c.Sort.Order == Desc {
			order = -1
		}
		pipeline = append(pipeline, bson.M{"$sort": bson.M{c.Sort.Field: order}})
	}

	// Apply pagination
	if c.Pagination.Offset > 0 {
		pipeline = append(pipeline, bson.M{"$skip": c.Pagination.Offset})
	}
	if c.Pagination.Limit > 0 {
		pipeline = append(pipeline, bson.M{"$limit": c.Pagination.Limit})
	}

	return pipeline
}
