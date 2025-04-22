package main

import (
	"fmt"
	"net/http"

	"github.com/kudras3r/KDTog/internal/chat/ws"
	"github.com/kudras3r/KDTog/internal/server"
	"github.com/kudras3r/KDTog/pkg/config"
	"github.com/kudras3r/KDTog/pkg/logger"
)

func main() {
	// init config
	config := config.Load()

	// init logger
	log := logger.New(config.LogLevel)
	log.Infof("log level set to %s", config.LogLevel)
	ws.SetLogger(log)

	// set ws config
	ws.SetConfig(config.WSock)
	log.Infof("ws upgrader set with read buffer size: %d, write buffer size: %d",
		config.WSock.RWBuffSize, config.WSock.RWBuffSize)
	log.Infof("max message size: %d", config.WSock.MaxMessSize)

	// init hub
	hub := ws.NewHub(log)
	log.Info("starting hub")
	go hub.Run()

	// init router
	r := server.NewRouter(hub, log)

	// run server
	log.Infof("starting server on %s:%s", config.Server.Host, config.Server.Port)
	if err := http.ListenAndServe(
		fmt.Sprintf("%s:%s", config.Server.Host, config.Server.Port),
		r); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
