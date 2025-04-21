package main

import (
	"log"
	"net/http"

	"github.com/kudras3r/KDTog/internal/chat/ws"
	"github.com/kudras3r/KDTog/internal/server"
)

func main() {
	// TODO
	// init config
	// init logger

	// init hub
	hub := ws.NewHub()
	go hub.Run()

	// init router
	r := server.NewRouter(hub)
	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
