package dbpool

import (
	"auth-service/pkg/config"
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
)

var maxRetries = 5

func BuildDSN(cfg *config.Config) string {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		cfg.DB.User,
		cfg.DB.Pass,
		cfg.DB.Host,
		cfg.DB.Port,
		cfg.DB.Name,
		cfg.DB.SSLMode,
	)

	return dsn
}

func InitDBPool(dsn string, log *logrus.Logger, ctx context.Context) (*pgxpool.Pool, error) {
	if dsn == "" {
		log.WithFields(logrus.Fields{
			"dsn": dsn,
		}).Error("empty dsn")
		return nil, fmt.Errorf("empty dsn")
	}

	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		log.WithError(err).WithFields(logrus.Fields{
			"dsn": dsn,
		}).Error("failed parse config")
		return nil, err
	}

	config.MaxConns = 10
	config.HealthCheckPeriod = 1 * time.Minute
	config.MaxConnLifetime = 1 * time.Hour
	config.MaxConnIdleTime = 5 * time.Second
	config.MinIdleConns = 2

	var pool *pgxpool.Pool
	for i := 0; i < maxRetries; i++ {
		pool, err = pgxpool.NewWithConfig(ctx, config)
		if err == nil {
			if err = pool.Ping(ctx); err == nil {
				log.Info("Pool initialize succeeded")
				return pool, nil
			}
			log.WithError(err).Errorf("database ping failed: %v", err)
		}
		log.WithError(err).WithFields(logrus.Fields{
			"config": config,
		}).Error("failed create new pool with config")
		time.Sleep(5 * time.Second)
	}

	defer pool.Close()

	log.Errorf("failed connect to Database: %v", err)
	return nil, err
}
