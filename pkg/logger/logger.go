package logger

import (
	"auth-service/pkg/config"
	"os"

	"github.com/sirupsen/logrus"
)

func LoggerInit(cfg *config.Config) *logrus.Logger {
	log := logrus.New()

	log.SetFormatter(&logrus.JSONFormatter{})

	level, err := logrus.ParseLevel(cfg.Logger.Level)
	if err != nil {
		level = logrus.InfoLevel
	}

	log.SetLevel(level)

	log.SetOutput(os.Stdout)

	log = log.WithFields(logrus.Fields{
		"service": "auth-service",
	}).Logger

	log.WithField("service", "auth-service").Info("Logger initialized")

	return log
}
