package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

func New() *logrus.Logger {
	file, err := os.OpenFile("main.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		logrus.Fatalf("Error occured while opening log file: %v", err)
	}

	logger := logrus.New()
	logger.SetLevel(logrus.InfoLevel)
	logger.SetLevel(logrus.ErrorLevel)
	logger.SetOutput(file)
	logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
		DisableColors: false,
	})

	return logger
}
