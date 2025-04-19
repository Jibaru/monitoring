package domain

import "time"

type Log struct {
	id        ID
	appID     ID
	timestamp time.Time
	data      map[string]interface{}
	raw       string
	level     string
}

func NewLog(
	id ID,
	appID ID,
	timestamp time.Time,
	data map[string]interface{},
	raw string,
	level string,
) (*Log, error) {
	return &Log{
		id:        id,
		appID:     appID,
		timestamp: timestamp,
		data:      data,
		raw:       raw,
		level:     level,
	}, nil
}
