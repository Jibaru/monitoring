package handler

import (
	"monitoring/config"
	"monitoring/internal/db"
	"monitoring/server"
	"net/http"
)

var cfg = config.Load()

func Handler(w http.ResponseWriter, r *http.Request) {
	db, client := db.New(cfg)
	defer client.Disconnect(r.Context())
	server.New(cfg, db).ServeHTTP(w, r)
}
