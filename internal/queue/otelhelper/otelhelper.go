package otelhelper

import (
	"context"
	"net/http"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"go.opentelemetry.io/otel/trace"
	"golang.org/x/exp/maps"
)

// StartProducerSpan abstracts the span creation for producing.
func StartProducerSpan(
	ctx context.Context,
	tracer trace.Tracer,
	name string,
) (context.Context, trace.Span) {
	attr := []attribute.KeyValue{
		semconv.MessagingSystem("memory"),
		semconv.MessagingOperationPublish,
	}

	ctx, span := tracer.Start(ctx, name,
		trace.WithAttributes(attr...),
		trace.WithSpanKind(trace.SpanKindProducer),
	)

	return ctx, span
}

// StartConsumerSpan abstracts span consumption
func StartConsumerSpan(
	tracer trace.Tracer,
	name string,
	headers http.Header,
) (context.Context, trace.Span) {
	attr := []attribute.KeyValue{
		semconv.MessagingSystem("memory"),
		semconv.MessagingOperationPublish,
	}

	propagator := otel.GetTextMapPropagator()

	ctx := propagator.Extract(context.Background(), carrier{headers})
	ctx, span := tracer.Start(ctx, name,
		trace.WithAttributes(attr...),
		trace.WithSpanKind(trace.SpanKindConsumer),
	)
	// propagator.Inject(ctx, carrier{headers})

	return ctx, span
}

type carrier struct{ h http.Header }

func (c carrier) Get(key string) string {
	return c.h.Get(key)
}

func (c carrier) Set(key, val string) {
	c.h.Set(key, val)
}

// Keys of the carrier
func (c carrier) Keys() []string {
	return maps.Keys(c.h)
}
