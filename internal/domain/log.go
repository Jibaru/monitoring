package domain

import (
	"encoding/json"
	"time"
)

type Log struct {
	id        ID
	appID     ID
	timestamp time.Time
	data      map[string]any
	raw       string
	level     string
}

func NewLog(
	id ID,
	appID ID,
	timestamp time.Time,
	data map[string]any,
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

func (l *Log) ID() ID {
	return l.id
}

func (l *Log) AppID() ID {
	return l.appID
}

func (l *Log) Timestamp() time.Time {
	return l.timestamp
}

func (l *Log) Data() map[string]any {
	return l.data
}

func (l *Log) Raw() string {
	return l.raw
}

func (l *Log) Level() string {
	return l.level
}

func (a Log) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]any{
		"id":        a.id,
		"appId":     a.appID,
		"timestamp": a.timestamp,
		"data":      a.data,
		"raw":       a.raw,
		"level":     a.level,
	})
}
