package config

import (
	"log"
	"os"
	"path/filepath"

	"github.com/codingconcepts/env"
	"github.com/joho/godotenv"
)

type Config struct {
	Server   Server `env:"SERVER"`
	WSock    WSock  `env:"WS"`
	LogLevel string `env:"LOG_LEVEL"`
}

type Server struct {
	Host string `env:"SERVER_HOST"`
	Port string `env:"SERVER_PORT"`
}

type WSock struct {
	RWBuffSize  int `env:"WS_RW_BUFF_SIZE"`
	MaxMessSize int `env:"WS_MAX_MESS_SIZE"`
}

func Load() *Config {
	envPath := filepath.Join(".", ".env")
	if err := godotenv.Load(envPath); err != nil {
		log.Fatal(err)
	}

	var config Config
	var serverConfig Server
	var wsockConfig WSock

	if err := env.Set(&wsockConfig); err != nil {
		log.Fatal("cannot get wsock env vars: ", err)
	}
	if err := env.Set(&serverConfig); err != nil {
		log.Fatal("cannot get server env vars: ", err)
	}
	config.LogLevel = os.Getenv("LOG_LEVEL")

	config.WSock = wsockConfig
	config.Server = serverConfig

	return &config
}
