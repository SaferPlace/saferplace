package saferplace

import (
	"context"
	"fmt"
	"log/slog"
	"slices"

	"golang.org/x/exp/maps"
	"golang.org/x/sync/errgroup"
	"safer.place/internal/config"
	"safer.place/internal/review"
	"safer.place/internal/service"

	// Registered services
	"safer.place/internal/service/imageupload"
	reportv1 "safer.place/internal/service/report/v1"
	reviewv1 "safer.place/internal/service/review/v1"
	viewerv1 "safer.place/internal/service/viewer/v1"
)

// Component is a type which contains all the component definitions
type Component string

const (
	ConsumerComponent Component = "consumer"
	ReviewComponent   Component = "review"
	ReportComponent   Component = "report"
	UploaderComponent Component = "uploader"
	ViewerComponent   Component = "viewer"
)

var componentDependencies = map[Component][]Dependency{
	ConsumerComponent: {QueueDependency, DatabaseDependency, NotifierDependency},
	ReviewComponent:   {DatabaseDependency},
	ReportComponent:   {QueueDependency},
	UploaderComponent: {StorageDependency},
	ViewerComponent:   {DatabaseDependency},
}

var headlessComponents = map[Component]registerHeadlessComponentFn{
	ConsumerComponent: registerConsumer,
}

type ComponentRegisterMap = map[Component]registerComponentFn

var reviewerComponents = ComponentRegisterMap{
	ReviewComponent: registerReview,
}

var userComponents = ComponentRegisterMap{
	ReportComponent:   registerReport,
	UploaderComponent: registerUploader,
	ViewerComponent:   registerViewer,
}

// StringsToComponents convert string slice to component slice or panic
// if an unknown component has been encountered.
func StringsToComponents(ss []string) []Component {
	res := make([]Component, 0, len(ss))
	for _, s := range ss {
		switch s {
		case string(ConsumerComponent):
			res = append(res, ConsumerComponent)
		case string(ReviewComponent):
			res = append(res, ReportComponent)
		case string(ReportComponent):
			res = append(res, ReportComponent)
		case string(UploaderComponent):
			res = append(res, UploaderComponent)
		case string(ViewerComponent):
			res = append(res, ViewerComponent)
		default:
			panic(fmt.Sprintf("unrecognised component %q", s))
		}
	}
	return res
}

// ComponentsToStrings returns a slice of strings from the slice of components.
func ComponentsToStrings(components []Component) []string {
	res := make([]string, 0, len(components))
	for _, component := range components {
		res = append(res, string(component))
	}
	return res
}

// AllComponents returns all the known components
func AllComponents() []Component {
	return maps.Keys(componentDependencies)
}

// neededDependencies returns a list of all dependencies that are needed for the given components.
func neededDependencies(components []Component) []Dependency {
	var dependencies []Dependency
	for _, component := range components {
		dependencies = append(dependencies, componentDependencies[component]...)
	}
	return dependencies
}

type registerHeadlessComponentFn func(context.Context, *config.Config, *dependencies, *errgroup.Group) error

func createHeadlessComponents(ctx context.Context, cfg *config.Config, wantedComponents []Component, deps *dependencies, eg *errgroup.Group) error {
	for component, fn := range headlessComponents {
		if slices.Contains(wantedComponents, component) {
			return fn(ctx, cfg, deps, eg)
		}
	}

	return nil
}

type registerComponentFn func(context.Context, *config.Config, *dependencies) (service.Service, error)

func createServices(ctx context.Context, cfg *config.Config, wantedComponents []Component, deps *dependencies, m ComponentRegisterMap) ([]service.Service, error) {
	services := make([]service.Service, 0, len(wantedComponents))
	for component, fn := range m {
		if slices.Contains(wantedComponents, component) {
			service, err := fn(ctx, cfg, deps)
			if err != nil {
				return nil, err
			}
			services = append(services, service)
		}
	}
	return services, nil
}

func registerConsumer(ctx context.Context, cfg *config.Config, deps *dependencies, eg *errgroup.Group) error {
	consumer := review.New(
		deps.logger.With(slog.String("component", "review")),
		deps.queue,
		deps.database,
		deps.notifer,
	)

	eg.Go(func() error {
		return consumer.Run(ctx)
	})

	return nil
}

func registerReview(_ context.Context, _ *config.Config, deps *dependencies) (service.Service, error) {
	return reviewv1.Register(
		reviewv1.Database(deps.database),
		reviewv1.Logger(deps.logger.With(slog.String("service", "reviewv1"))),
		reviewv1.Tracer(deps.tracing.Tracer("review")),
	), nil
}

func registerReport(_ context.Context, _ *config.Config, deps *dependencies) (service.Service, error) {
	return reportv1.Register(
		deps.queue,
		deps.logger.With(slog.String("service", "reportv1")),
	), nil
}

func registerUploader(_ context.Context, _ *config.Config, deps *dependencies) (service.Service, error) {
	return imageupload.Register(
		imageupload.Logger(deps.logger.With(slog.String("service", "imageupload"))),
		imageupload.Tracer(deps.tracing.Tracer("imageupload")),
		imageupload.Storage(deps.storage),
	), nil
}

func registerViewer(_ context.Context, _ *config.Config, deps *dependencies) (service.Service, error) {
	return viewerv1.Register(
		deps.database,
		deps.logger.With(slog.String("service", "viewerv1")),
	), nil
}
