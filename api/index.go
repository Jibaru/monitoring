package handler

import (
	"monitoring/config"
	"monitoring/server"
	"net/http"
)

var cfg = config.Load()
var mux = server.New(cfg)

func Handler(w http.ResponseWriter, r *http.Request) {
	mux.ServeHTTP(w, r)
}
