package server

import (
	"net/http"

	"github.com/kudras3r/KDTog/internal/chat/ws"
	"github.com/kudras3r/KDTog/pkg/logger"
)

func getIndexHandler(log *logger.Logger) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Info("index handler called")
		http.ServeFile(w, r, "static/index.html")
	}
}

func getWSHandler(hub *ws.Hub) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ws.ServeWs(hub, w, r)
	}
}
