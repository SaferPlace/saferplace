package saferplace

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"go.uber.org/zap"
	"golang.org/x/exp/slices"
	"golang.org/x/sync/errgroup"
	"safer.place/internal/auth"
	"safer.place/internal/config"
	"safer.place/internal/review"
	"safer.place/webserver"
	"safer.place/webserver/certificate"
	"safer.place/webserver/certificate/insecure"
	"safer.place/webserver/certificate/temporary"
	"safer.place/webserver/middleware"

	// Registered services
	"safer.place/internal/service/imageupload"
	reportv1 "safer.place/internal/service/report/v1"
	reviewv1 "safer.place/internal/service/review/v1"
	viewerv1 "safer.place/internal/service/viewer/v1"
)

func Run(components []string, cfg *config.Config) (err error) {
	logger := newLogger(cfg)
	defer func() { _ = logger.Sync() }()

	logger.Debug("using components",
		zap.Strings("components", components),
	)

	eg, ctx := errgroup.WithContext(context.Background())

	// Setup all dependencies
	db, err := newDatabase(cfg.Database)
	if err != nil {
		return err
	}

	incidentQueue, err := newQueue(cfg.Queue)
	if err != nil {
		return err
	}

	store, err := newStorage(cfg.Storage)
	if err != nil {
		return err
	}

	notifier, err := newNotififer(cfg.Notifier, logger)
	if err != nil {
		return err
	}

	// Setup Components
	if slices.Contains(components, "consumer") {
		consumer := review.New(
			logger.With(zap.String("component", "review")),
			incidentQueue,
			db,
			notifier,
		)

		eg.Go(func() error {
			return consumer.Run(ctx)
		})
	}

	reviewerServices := []webserver.Service{}

	if slices.Contains(components, "review") {
		reviewerServices = append(reviewerServices,
			reviewv1.Register(
				db,
				logger.With(zap.String("service", "reviewv1")),
				// TODO: Re-enable once we know what we are doing.
				// auth.NewAuthInterceptor(db),
			),
		)
	}

	userServices := []webserver.Service{}

	if slices.Contains(components, "report") {
		userServices = append(userServices,
			reportv1.Register(incidentQueue, logger.With(zap.String("service", "reportv1"))),
		)
	}

	if slices.Contains(components, "uploader") {
		userServices = append(userServices,
			imageupload.Register(logger.With(zap.String("service", "imageupload")), store),
		)
	}

	if slices.Contains(components, "viewer") {
		userServices = append(userServices,
			viewerv1.Register(db, logger.With(zap.String("service", "viewerv1"))),
		)
	}

	// Setup Webserver based on the provided services
	userAuthMiddleware := auth.NewUserAuthMiddleware()
	services := append(
		reviewerServices,
		ServiceMiddleware(
			[]middleware.Middleware{userAuthMiddleware},
			userServices,
		)...,
	)

	middlewares := []middleware.Middleware{
		middleware.Cors(cfg.Webserver.CORSDomains),
	}

	var certProvider certificate.Provider
	switch cfg.Webserver.Cert.Provider {
	case "temporary":
		certProvider = temporary.NewProvider(temporary.Config{
			ValidFor: time.Hour,
		})
	case "insecure":
		certProvider = insecure.NewProvider()
	}

	tlsConfig, err := certProvider.Provide(context.Background(), cfg.Webserver.Cert.Domains)
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

	eg.Go(func() error {
		return srv.Run(cfg.Webserver.Port)
	})

	return eg.Wait()
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
