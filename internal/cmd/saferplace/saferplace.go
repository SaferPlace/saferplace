package saferplace

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	_ "net/http/pprof"

	"connectrpc.com/connect"
	"connectrpc.com/otelconnect"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/saferplace/webserver-go"
	"github.com/saferplace/webserver-go/middleware"
	"golang.org/x/sync/errgroup"
	"safer.place/internal/auth"
	"safer.place/internal/config"
	"safer.place/internal/service"
)

func Run(ctx context.Context, components []Component, cfg *config.Config) (err error) {
	eg, ctx := errgroup.WithContext(ctx)

	// Setup all deps
	deps, depCloser, err := createDependencies(ctx, cfg, components)
	if err != nil {
		return
	}
	defer depCloser.Close()

	// shared middleware
	middlewares := []middleware.Middleware{
		middleware.Cors(cfg.Webserver.CORSDomains),
	}

	// shared interceptors
	tracingInteceptor, err := otelconnect.NewInterceptor(
		otelconnect.WithTracerProvider(deps.tracing),
	)
	if err != nil {
		return fmt.Errorf("unable to setup tracing interceptor: %w", err)
	}
	interceptors := []connect.Interceptor{
		tracingInteceptor,
	}

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

	// creates services with the internal services
	services := []webserver.Service{
		profile,
		metrics(deps.metrics),
	}
	services = append(services,
		FinalizeServices(
			nil,
			interceptors,
			reviewerServices,
		)...,
	)
	services = append(services,
		FinalizeServices(
			[]middleware.Middleware{userAuthMiddleware},
			interceptors,
			userServices,
		)...,
	)

	tlsConfig, err := newTLSConfig(ctx, cfg.Webserver.Cert)
	if err != nil {
		return err
	}

	srv, err := webserver.New(
		webserver.Logger(deps.logger.Unwrap().With(slog.String("component", "server"))),
		webserver.Services(services...),
		webserver.TLSConfig(tlsConfig),
		webserver.Middlewares(middlewares...),
		webserver.ReadTimeout(cfg.Webserver.ReadTimeout),
		webserver.WriteTimeout(cfg.Webserver.WriteTimeout),
	)
	if err != nil {
		return fmt.Errorf("unable to create the server: %w", err)
	}

	eg.Go(func() error {
		return srv.Run(cfg.Webserver.Port)
	})

	return eg.Wait()
}

// FinalizeServices wraps all provided services with the middleware.
func FinalizeServices(
	middlewares []middleware.Middleware,
	interceptors []connect.Interceptor,
	services []service.Service,
) []webserver.Service {
	wrapped := make([]webserver.Service, 0, len(services))

	for _, service := range services {
		path, handler := service(interceptors...)
		for _, middleware := range middlewares {
			handler = middleware(handler)
		}
		wrapped = append(wrapped, func() (string, http.Handler) {
			return path, handler
		})
	}

	return wrapped
}

func metrics(reg *prometheus.Registry) func() (string, http.Handler) {
	return func() (string, http.Handler) {
		return "/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{
			EnableOpenMetrics: true,
		})
	}
}

func profile() (string, http.Handler) {
	// Get the default mux as pprof registers correctly with it
	mux := http.DefaultServeMux

	return "/debug/pprof/", mux
}
