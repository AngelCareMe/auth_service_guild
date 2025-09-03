package logger

import (
	"auth-service/pkg/config"
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
)

func LoggerInit(cfg *config.Config) (*logrus.Logger, error) {
	log := logrus.New()

	log.SetFormatter(&logrus.JSONFormatter{})

	level, err := logrus.ParseLevel(cfg.Logger.Level)
	if err != nil {
		return nil, fmt.Errorf("failed parse log level: %w", err)
	}

	log.SetLevel(level)

	log.SetOutput(os.Stdout)

	log = log.WithFields(logrus.Fields{
		"service": "auth-service",
	}).Logger

	return log, nil
}
