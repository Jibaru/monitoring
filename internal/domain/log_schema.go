package domain

type LogSchema struct {
	Total  int            `json:"total"`
	Schema map[string]int `json:"schema"`
}
