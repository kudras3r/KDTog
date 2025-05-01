package server

import (
	"net/http"

	"github.com/kudras3r/KDTog/internal/chat/ws"
	"github.com/kudras3r/KDTog/internal/storage"
	"github.com/kudras3r/KDTog/pkg/logger"
)

const (
	GLOC = "internal/server/"
)

func NewRouter(log *logger.Logger, hub *ws.Hub, db storage.Storage) http.Handler {
	mux := http.NewServeMux()

	mux.Handle("/", authMiddleware(http.HandlerFunc(indexHandler(log)), log))
	// mux.Handle("/", http.HandlerFunc(indexHandler(log)))
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	mux.Handle("/auth", http.HandlerFunc(authHandler(log, db)))
	mux.Handle("/ws", http.HandlerFunc(WSHandler(hub)))

	return mux
}
