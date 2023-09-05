package saferplace

import (
	"context"
	"fmt"
	"net/http"

	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"safer.place/internal/auth"
	"safer.place/internal/config"
	"safer.place/webserver"
	"safer.place/webserver/middleware"
)

// Service is webserver registered function to create a new service, aliased for convenience
type Service = webserver.Service

func Run(components []Component, cfg *config.Config) (err error) {
	// Setup all deps
	deps, err := createDependencies(cfg, components)
	if err != nil {
		return
	}
	defer func() { _ = deps.logger.Sync() }()

	eg, ctx := errgroup.WithContext(context.Background())

	if err := createHeadlessComponents(ctx, cfg, components, deps, eg); err != nil {
		return fmt.Errorf("unable to create headless components: %w", err)
	}

	reviewerServices, err := createServices(ctx, cfg, components, deps, reviewerComponents)
	if err != nil {
		return fmt.Errorf("unable to create reviewer services: %w", err)
	}

	userServices, err := createServices(ctx, cfg, components, deps, userComponents)
	if err != nil {
		return fmt.Errorf("unable to create user services: %w", err)
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

	tlsConfig, err := newTLSConfig(cfg.Webserver.Cert)
	if err != nil {
		return err
	}

	srv, err := webserver.New(
		webserver.Logger(deps.logger.With(zap.String("component", "server"))),
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
