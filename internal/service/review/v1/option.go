package review

import (
	"go.opentelemetry.io/otel/trace"
	"safer.place/internal/database"
	"safer.place/internal/log"
)

type Option func(*Service)

func Logger(l log.Logger) Option {
	return func(s *Service) {
		s.log = l
	}
}

func Tracer(t trace.Tracer) Option {
	return func(s *Service) {
		s.tracer = t
	}
}

func Database(db database.Review) Option {
	return func(s *Service) {
		s.db = db
	}
}
