package saferplace

import (
	"errors"
	"fmt"

	"api.safer.place/incident/v1"
	"go.uber.org/zap"
	"safer.place/internal/config"
	"safer.place/internal/database"
	"safer.place/internal/database/sqldatabase"
	"safer.place/internal/notifier"
	"safer.place/internal/notifier/lognotifier"
	"safer.place/internal/queue"
	"safer.place/internal/queue/memory"
	"safer.place/internal/storage"
	"safer.place/internal/storage/minio"
)

var errProviderNotFound = errors.New("provider not found")

func newLogger(cfg *config.Config) *zap.Logger {
	var logger *zap.Logger
	if cfg.Debug {
		logger, _ = zap.NewDevelopment()
		logger.Debug("debug mode enabled")
	} else {
		logger, _ = zap.NewProduction()
	}

	logger.Debug("using configuration",
		zap.Any("config", cfg),
	)

	return logger
}

func newDatabase(cfg config.DatabaseConfig) (v database.Database, err error) {
	switch cfg.Provider {
	case "sql":
		v, err = sqldatabase.New(cfg.SQL)
	default:
		err = errProviderNotFound
	}

	if err != nil {
		return nil, fmt.Errorf("unable to open %q database: %w", cfg.Provider, err)
	}

	return v, nil
}

func newQueue(cfg config.QueueConfig) (v queue.Queue[*incident.Incident], err error) {
	switch cfg.Provider {
	case "memory":
		v = memory.New[*incident.Incident]()
	default:
		err = errProviderNotFound
	}

	if err != nil {
		return nil, fmt.Errorf("unable to open %q queue: %w", cfg.Provider, err)
	}

	return v, nil
}

func newStorage(cfg config.StorageConfig) (v storage.Storage, err error) {
	switch cfg.Provider {
	case "minio":
		v, err = minio.New(cfg.Minio)
	default:
		err = errProviderNotFound
	}

	if err != nil {
		return nil, fmt.Errorf("unable to open %q storage: %w", cfg.Provider, err)
	}

	return v, nil
}

func newNotififer(cfg config.NotifierConfig, logger *zap.Logger) (v notifier.Notifier, err error) {
	log := logger.With(zap.String("notifier", cfg.Provider))
	switch cfg.Provider {
	case "log":
		v = lognotifier.New(log)
	default:
		err = errProviderNotFound
	}

	if err != nil {
		return nil, fmt.Errorf("unable to open %q database: %w", cfg.Provider, err)
	}

	return v, nil
}
