package minio

import (
	"go.opentelemetry.io/otel/trace"
)

// Option extends the functionality of the storage
type Option func(*Storage)

// Tracer provides the tracing
func Tracer(t trace.Tracer) Option {
	return func(s *Storage) {
		s.tracer = t
	}
}
