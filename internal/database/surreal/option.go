package surreal

import (
	"go.opentelemetry.io/otel/trace"
	"safer.place/internal/log"
)

type Option func(*Database)

func Logger(l log.Logger) Option {
	return func(db *Database) {
		db.logger = l
	}
}

func Tracer(t trace.Tracer) Option {
	return func(db *Database) {
		db.tracer = t
	}
}
