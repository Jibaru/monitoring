package domain

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

var (
	ErrApp = fmt.Errorf("error in app")
)

type App struct {
	id        ID
	name      string
	appKey    string
	userID    ID
	createdAt time.Time
}

func NewApp(
	id ID,
	name string,
	appKey string,
	userID ID,
	createdAt time.Time,
) (*App, error) {
	app := &App{
		id:        id,
		userID:    userID,
		createdAt: createdAt,
	}

	if err := app.ChangeName(name); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrApp, err)
	}

	if err := app.ChangeAppKey(appKey); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrApp, err)
	}
	return app, nil
}

func (a *App) ID() ID {
	return a.id
}

func (a *App) Name() string {
	return a.name
}

func (a *App) AppKey() string {
	return a.appKey
}

func (a *App) UserID() ID {
	return a.userID
}

func (a *App) CreatedAt() time.Time {
	return a.createdAt
}

func (a *App) ChangeName(name string) error {
	if strings.TrimSpace(name) == "" {
		return fmt.Errorf("%w: name cannot be empty", ErrApp)
	}

	a.name = name
	return nil
}

func (a *App) ChangeAppKey(appKey string) error {
	if strings.TrimSpace(appKey) == "" {
		return fmt.Errorf("%w: appKey cannot be empty", ErrApp)
	}

	a.appKey = appKey
	return nil
}

func (a App) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]any{
		"id":        a.id,
		"name":      a.name,
		"appKey":    a.appKey,
		"userId":    a.userID,
		"createdAt": a.createdAt,
	})
}
