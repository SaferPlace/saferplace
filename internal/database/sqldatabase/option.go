package sqldatabase

import (
	"errors"

	"go.opentelemetry.io/otel/trace"
)

type Option func(db *Database)

func Tracer(tracer trace.Tracer) Option {
	return func(db *Database) {
		db.tracer = tracer
	}
}

var (
	errMissingTracer = errors.New("missing tracer")
)

func validate(db *Database) error {
	if db.tracer == nil {
		return errMissingTracer
	}

	return nil
}
