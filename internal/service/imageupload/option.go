package imageupload

import (
	"go.opentelemetry.io/otel/trace"

	"safer.place/internal/log"
	"safer.place/internal/storage"
)

// Option to provide configuration to the service.
type Option func(*Service)

// Logger provides the logger
func Logger(log log.Logger) Option {
	return func(s *Service) {
		s.log = log
	}
}

// Trace provides the tracing
func Tracer(tp trace.Tracer) Option {
	return func(s *Service) {
		s.tracer = tp
	}
}

// Storage provides the storage
func Storage(store storage.Storage) Option {
	return func(s *Service) {
		s.storage = store
	}
}
