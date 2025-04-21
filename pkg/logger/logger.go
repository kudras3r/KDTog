package logger

import "github.com/sirupsen/logrus"

type Logger struct {
	logrus.Logger
}

func New() *Logger {
	logger := &Logger{
		Logger: *logrus.New(),
	}

	logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
		DisableColors:   false,
	})
	logger.SetLevel(logrus.DebugLevel)

	return logger
}
