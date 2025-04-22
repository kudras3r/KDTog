package logger

import (
	"strings"

	"github.com/sirupsen/logrus"
)

const (
	DEDUG = "debug"
	INFO  = "info"
)

type Logger struct {
	logrus.Logger
}

func New(level string) *Logger {
	logger := &Logger{
		Logger: *logrus.New(),
	}

	logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
		DisableColors:   false,
	})
	switch strings.ToLower(level) {
	case INFO:
		logger.SetLevel(logrus.InfoLevel)
	case DEDUG:
		logger.SetLevel(logrus.DebugLevel)
	default:
		logger.SetLevel(logrus.InfoLevel)
	}
	return logger
}
