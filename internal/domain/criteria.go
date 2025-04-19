package domain

// FilterType represents the type of filter (equals or not equals).
type FilterType string

const (
	Equals             FilterType = "eq"
	NotEquals          FilterType = "ne"
	Like               FilterType = "like"
	In                 FilterType = "in"
	GreaterThanOrEqual FilterType = "gte"
	LessThanOrEqual    FilterType = "lte"
)

var (
	EmptyPagination Pagination = Pagination{}
	EmptySort       Sort       = Sort{}
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
