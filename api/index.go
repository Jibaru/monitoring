package handler

import (
	"net/http"

	"monitoring/config"
	"monitoring/db"
	"monitoring/server"
)

var cfg = config.Load()

func Handler(w http.ResponseWriter, r *http.Request) {
	db, client := db.New(cfg)
	defer client.Disconnect(r.Context())
	server.New(cfg, db).ServeHTTP(w, r)
}
