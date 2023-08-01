package realtime

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"safer.place/realtime/internal/config"
	"safer.place/realtime/internal/database"
	"safer.place/realtime/internal/database/sqldatabase"
	"safer.place/realtime/internal/notifier"
	"safer.place/realtime/internal/notifier/discordnotifier"
	"safer.place/realtime/internal/notifier/lognotifier"
	"safer.place/realtime/internal/queue"
	"safer.place/realtime/internal/queue/memory"
	"safer.place/realtime/internal/review"
	reviewui "safer.place/realtime/packages/review-ui"
	"safer.place/webserver"
	"safer.place/webserver/certificate"
	"safer.place/webserver/certificate/insecure"
	"safer.place/webserver/certificate/temporary"

	"api.safer.place/incident/v1"

	// Registered services
	reportv1 "safer.place/realtime/internal/service/report/v1"
	reviewv1 "safer.place/realtime/internal/service/review/v1"
	"safer.place/realtime/internal/service/static"
)

func Run(cfg *config.Config) (err error) {
	var logger *zap.Logger
	if cfg.Debug {
		logger, _ = zap.NewDevelopment()
		logger.Debug("debug mode enabled")
	} else {
		logger, _ = zap.NewProduction()
	}
	defer func() { _ = logger.Sync() }()

	var q queue.Queue[*incident.Incident]
	switch cfg.Queue {
	case "memory":
		q = memory.New[*incident.Incident]()
	}

	var db database.Database
	switch cfg.Database {
	case "sql":
		db, err = sqldatabase.New()
		if err != nil {
			return fmt.Errorf("unable to open SQL database: %w", err)
		}
	}

	var dn notifier.Notifier
	switch cfg.Notifier {
	case "discord":
		dn, err = discordnotifier.New(http.DefaultClient)
		if err != nil {
			return fmt.Errorf("unable to create discord notifier: %w", err)
		}
	case "log":
		dn = lognotifier.New(logger.With(zap.String("notifier", "log")))
	default:
		return fmt.Errorf("notifier not specified")
	}

	r := review.New(
		logger.With(zap.String("component", "review")),
		q,
		db,
		dn,
	)

	services := []webserver.Service{
		reportv1.Register(q, logger.With(zap.String("service", "reportv1"))),
		reviewv1.Register(db, logger.With(zap.String("service", "reviewv1"))),
		// TODO: Once we add more frontends maybe it would be better to move
		// somewhere better.
		static.Register("/", reviewui.StaticFiles), // review UI
	}

	var certProvider certificate.Provider
	switch cfg.Cert.Provider {
	case "temporary":
		certProvider = temporary.NewProvider(temporary.Config{
			ValidFor: time.Hour,
		})
	case "insecure":
		certProvider = insecure.NewProvider()
	}

	tlsConfig, err := certProvider.Provide(context.Background(), cfg.Cert.Domains)
	if err != nil {
		return fmt.Errorf("unable to create TLS cert: %w", err)
	}

	srv, err := webserver.New(
		webserver.Logger(logger.With(zap.String("component", "server"))),
		webserver.Services(services...),
		webserver.TLSConfig(tlsConfig),
	)
	if err != nil {
		return fmt.Errorf("unable to create the server: %w", err)
	}

	eg, ctx := errgroup.WithContext(context.Background())
	eg.Go(func() error {
		return r.Run(ctx)
	})
	eg.Go(func() error {
		return srv.Run(cfg.Port)
	})

	return eg.Wait()
}
