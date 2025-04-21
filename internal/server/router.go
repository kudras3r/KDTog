package server

import (
	"net/http"

	"github.com/kudras3r/KDTog/internal/chat/ws"
)

func NewRouter(hub *ws.Hub) http.Handler {
	mux := http.NewServeMux()

	mux.Handle("/", http.HandlerFunc(indexHandler))
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	mux.Handle("/ws", http.HandlerFunc(getWSHandler(hub)))

	return mux
}
