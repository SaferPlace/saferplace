package tracing

import (
	"context"
	"fmt"
	"io"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
)

// Configuration used for tracing
type Config struct {
	Enabled bool `yaml:"enabled" default:"false"`

	Timeout       time.Duration `yaml:"timeout" default:"1s"`
	Endpoint      string        `yaml:"endpoint" default:"otel-collector:4317"`
	SamplingRatio float64       `yaml:"sampling_ratio" default:"1" split_words:"true"`
}

// NewTracingProvider creates the provider configuration from the config. If the
func NewTracingProvider(ctx context.Context, cfg *Config) (trace.TracerProvider, io.Closer, error) {
	if !cfg.Enabled {
		return trace.NewNoopTracerProvider(), noopCloser, nil
	}

	r, err := resource.New(ctx,
		resource.WithProcess(),
		resource.WithOS(),
		resource.WithContainer(),
		resource.WithHost(),
		resource.WithAttributes(semconv.ServiceName("saferplace")),
	)
	if err != nil {
		return nil, nil, fmt.Errorf("unable to ")
	}

	exporter, err := newExporter(ctx, cfg)
	if err != nil {
		return nil, nil, fmt.Errorf("unable to create a span processor: %w", err)
	}

	// Always sample unless the ratio is set lower.
	sampler := sdktrace.AlwaysSample()
	if cfg.SamplingRatio < 1 {
		sampler = sdktrace.ParentBased(
			sdktrace.TraceIDRatioBased(cfg.SamplingRatio),
		)
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSpanProcessor(exporter),
		sdktrace.WithResource(r),
		sdktrace.WithSampler(sampler),
	)

	otel.SetTracerProvider(tp)

	return tp, shutdownCloser(tp.Shutdown), nil
}

func newExporter(ctx context.Context, cfg *Config) (sdktrace.SpanProcessor, error) {
	ctx, cancel := context.WithTimeout(ctx, cfg.Timeout)
	defer cancel()

	client := otlptracegrpc.NewClient(
		otlptracegrpc.WithInsecure(),
		otlptracegrpc.WithEndpoint(cfg.Endpoint),
		otlptracegrpc.WithDialOption(grpc.WithBlock()))
	exporter, err := otlptrace.New(ctx, client)
	if err != nil {
		return nil, fmt.Errorf("unable to create client: %w", err)
	}

	return sdktrace.NewBatchSpanProcessor(exporter), nil
}

var noopCloser = closer(func() error { return nil })

type closer func() error

func (c closer) Close() error {
	return c()
}

type shutdownCloser func(context.Context) error

func (c shutdownCloser) Close() error {
	return c(context.Background())
}
