package realtime

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"safer.place/internal/auth"
	"safer.place/internal/config"
	"safer.place/internal/database"
	"safer.place/internal/database/sqldatabase"
	"safer.place/internal/notifier"
	"safer.place/internal/notifier/discordnotifier"
	"safer.place/internal/notifier/lognotifier"
	"safer.place/internal/queue"
	"safer.place/internal/queue/memory"
	"safer.place/internal/review"
	"safer.place/internal/storage"
	"safer.place/internal/storage/minio"
	reviewui "safer.place/packages/review-ui"
	"safer.place/webserver"
	"safer.place/webserver/certificate"
	"safer.place/webserver/certificate/insecure"
	"safer.place/webserver/certificate/temporary"
	"safer.place/webserver/middleware"

	"api.safer.place/incident/v1"

	// Registered services
	"safer.place/internal/service/imageupload"
	reportv1 "safer.place/internal/service/report/v1"
	reviewv1 "safer.place/internal/service/review/v1"
	viewerv1 "safer.place/internal/service/viewer/v1"
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

	logger.Debug("using configuration",
		zap.Any("config", cfg),
	)

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

	var store storage.Storage
	switch cfg.Storage {
	case "minio":
		store, err = minio.New()
		if err != nil {
			return fmt.Errorf("unable to open minio storage: %w", err)
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

	userAuthMiddleware := auth.NewUserAuthMiddleware()

	services := []webserver.Service{}
	// Reviewer services
	services = append(services,
		// Review Services
		reviewv1.Register(
			db,
			logger.With(zap.String("service", "reviewv1")),
			// TODO: Re-enable once we know what we are doing.
			// auth.NewAuthInterceptor(db),
		),

		// TODO: Once we add more frontends maybe it would be better to move
		// somewhere better.
		auth.Register("/review/", &auth.Config{
			Handler:      http.FileServer(http.FS(reviewui.StaticFiles)),
			Log:          logger.With(zap.String("service", "reviewauth")),
			Domain:       cfg.Auth.Domain,
			ClientID:     cfg.Auth.ClientID,
			ClientSecret: cfg.Auth.ClientSecret,
			DB:           db,
		}),
	)
	// User services
	services = append(services, ServiceMiddleware(
		[]middleware.Middleware{userAuthMiddleware},
		[]webserver.Service{
			reportv1.Register(q, logger.With(zap.String("service", "reportv1"))),
			viewerv1.Register(db, logger.With(zap.String("service", "viewerv1"))),
			imageupload.Register(logger.With(zap.String("service", "imageupload")), store),
		},
	)...)

	middlewares := []middleware.Middleware{
		middleware.Cors(nil),
	}
	// Only enable request logging when debug mode is active to prevent unnecessary collection.
	// This will be normally done using metrics.
	if cfg.Debug {
		middlewares = append(middlewares,
			loggingMiddleware(logger.With(zap.String("component", "requestlog"))),
		)
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
		webserver.Middlewares(middlewares...),
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

// TODO: Move this somewhere else
func loggingMiddleware(log *zap.Logger) middleware.Middleware {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Info("request",
				zap.String("path", r.URL.Path),
			)
			h.ServeHTTP(w, r)
		})
	}
}

// ServiceMiddleware wraps all provided services with the middleware.
func ServiceMiddleware(
	middlewares []middleware.Middleware, services []webserver.Service,
) []webserver.Service {
	wrapped := make([]webserver.Service, 0, len(services))

	for _, service := range services {
		path, handler := service()
		for _, middleware := range middlewares {
			handler = middleware(handler)
		}
		wrapped = append(wrapped, func() (string, http.Handler) {
			return path, handler
		})
	}

	return wrapped
}
