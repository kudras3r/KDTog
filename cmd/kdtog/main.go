package main

import (
	"net/http"

	"github.com/kudras3r/KDTog/internal/chat/ws"
	"github.com/kudras3r/KDTog/internal/server"
	"github.com/kudras3r/KDTog/pkg/logger"
)

func main() {
	// TODO
	// init config
	// init logger
	log := logger.New()
	ws.SetLogger(log)

	// init hub
	hub := ws.NewHub(log)
	log.Info("starting hub")
	go hub.Run()

	// init router
	r := server.NewRouter(hub, log)

	// run server
	log.Info("starting server on :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
