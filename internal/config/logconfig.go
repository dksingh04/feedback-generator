package config

import (
	"os"

	"github.com/sirupsen/logrus"
)

var logger *logrus.Logger

//CreateDefaultLogConfiguration used to create log configuration, which will be used in all over the application
func CreateDefaultLogConfiguration() (*logrus.Logger, error) {
	logger = logrus.New()
	logger.SetLevel(logrus.InfoLevel)
	logger.SetOutput(os.Stdout)
	logger.SetFormatter(&logrus.JSONFormatter{})
	return logger, nil
}

// GetDefaultLogger will return created Logger, in case you don't want to use this Logger you can create your own
func GetDefaultLogger() *logrus.Logger {
	return logger
}
