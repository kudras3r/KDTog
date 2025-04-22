package server

import (
	"net/http"

	"github.com/kudras3r/KDTog/internal/chat/ws"
	"github.com/kudras3r/KDTog/pkg/logger"
)

func NewRouter(hub *ws.Hub, log *logger.Logger) http.Handler {
	mux := http.NewServeMux()

	mux.Handle("/", http.HandlerFunc(getIndexHandler(log)))
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	mux.Handle("/ws", http.HandlerFunc(getWSHandler(hub)))

	return mux
}
