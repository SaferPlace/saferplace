package saferplace

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"slices"
	"time"

	"api.safer.place/incident/v1"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
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
	"safer.place/webserver/certificate"
	"safer.place/webserver/certificate/insecure"
	"safer.place/webserver/certificate/temporary"
)

var errProviderNotFound = errors.New("provider not found")

type Dependency string

const (
	DatabaseDependency Dependency = "database"
	QueueDependency    Dependency = "queue"
	StorageDependency  Dependency = "storage"
	NotifierDependency Dependency = "notifier"
)

func dependenciesToStrings(dependencies []Dependency) []string {
	res := make([]string, 0, len(dependencies))
	for _, dependency := range dependencies {
		res = append(res, string(dependency))
	}
	return res
}

// StringsToDependencies converts a string slice into dependecy slice
func StringsToDependencies(ss []string) []Dependency {
	res := make([]Dependency, 0, len(ss))
	for _, s := range ss {
		switch s {
		case string(DatabaseDependency):
			res = append(res, DatabaseDependency)
		case string(QueueDependency):
			res = append(res, QueueDependency)
		case string(StorageDependency):
			res = append(res, StorageDependency)
		case string(NotifierDependency):
			res = append(res, NotifierDependency)
		default:
			panic(fmt.Sprintf("unrecognised dependency %q", s))
		}
	}
	return res
}

type dependencies struct {
	metrics  *prometheus.Registry
	logger   *zap.Logger
	database database.Database
	queue    queue.Queue[*incident.Incident]
	storage  storage.Storage
	notifer  notifier.Notifier
}

type registerDependencyFn func(*config.Config, *dependencies) error

func createDependencies(cfg *config.Config, components []Component) (*dependencies, error) {
	wantedDependencies := neededDependencies(components)

	deps := &dependencies{
		logger:  newLogger(cfg),
		metrics: prometheus.NewRegistry(),
	}

	deps.logger.Debug("initializing dependencies",
		zap.Strings("components", ComponentsToStrings(components)),
		zap.Strings("dependencies", dependenciesToStrings(wantedDependencies)),
	)

	deps.metrics.MustRegister(
		collectors.NewGoCollector(),
		collectors.NewBuildInfoCollector(),
		collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
	)

	for dep, fn := range map[Dependency]registerDependencyFn{
		DatabaseDependency: registerDatabase,
		QueueDependency:    registerQueue,
		StorageDependency:  registerStorage,
		NotifierDependency: registerNotifier,
	} {
		if slices.Contains(wantedDependencies, dep) {
			if err := fn(cfg, deps); err != nil {
				return deps, err
			}
		}
	}

	return deps, nil
}

func newTLSConfig(cfg config.CertConfig) (v *tls.Config, err error) {
	var p certificate.Provider
	switch cfg.Provider {
	case "temporary":
		p = temporary.NewProvider(temporary.Config{
			ValidFor: time.Duration(cfg.ValidFor),
		})
	case "insecure":
		p = insecure.NewProvider()
	default:
		return nil, errProviderNotFound
	}

	v, err = p.Provide(context.Background(), cfg.Domains)
	if err != nil {
		return nil, fmt.Errorf("unable to create %q TLS config: %w", cfg.Provider, err)
	}

	return v, nil
}

func registerDatabase(cfg *config.Config, deps *dependencies) (err error) {
	var v database.Database
	switch cfg.Database.Provider {
	case "sql":
		v, err = sqldatabase.New(cfg.Database.SQL)
	default:
		err = errProviderNotFound
	}

	if err != nil {
		return fmt.Errorf("unable to open %q database: %w", cfg.Database.Provider, err)
	}

	deps.database = v
	return nil
}

func registerQueue(cfg *config.Config, deps *dependencies) (err error) {
	var v queue.Queue[*incident.Incident]
	switch cfg.Queue.Provider {
	case "memory":
		v = memory.New[*incident.Incident]()
	default:
		err = errProviderNotFound
	}

	if err != nil {
		return fmt.Errorf("unable to open %q queue: %w", cfg.Queue.Provider, err)
	}

	deps.queue = v
	return nil
}

func registerStorage(cfg *config.Config, deps *dependencies) (err error) {
	var v storage.Storage
	switch cfg.Storage.Provider {
	case "minio":
		v, err = minio.New(cfg.Storage.Minio)
	default:
		err = errProviderNotFound
	}

	if err != nil {
		return fmt.Errorf("unable to open %q storage: %w", cfg.Storage.Provider, err)
	}

	deps.storage = v
	return nil
}

func registerNotifier(cfg *config.Config, deps *dependencies) (err error) {
	var v notifier.Notifier
	log := deps.logger.With(zap.String("notifier", cfg.Notifier.Provider))
	switch cfg.Notifier.Provider {
	case "log":
		v = lognotifier.New(log)
	default:
		err = errProviderNotFound
	}

	if err != nil {
		return fmt.Errorf("unable to open %q database: %w", cfg.Notifier.Provider, err)
	}

	deps.notifer = v
	return nil
}

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