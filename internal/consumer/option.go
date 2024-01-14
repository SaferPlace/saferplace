package consumer

import (
	"go.opentelemetry.io/otel/trace"

	"safer.place/internal/database"
	"safer.place/internal/log"
	"safer.place/internal/notifier"
	"safer.place/internal/queue"

	"api.safer.place/incident/v1"
)

type Option func(r *Review)

// Consumer Option allows to specofy the consumer queue that we can get the messages from
func Consumer(c queue.Consumer[*incident.Incident]) Option {
	return func(r *Review) {
		r.incoming = c
	}
}

// Notifier is used to notify about incoming review
func Notifier(n notifier.Notifier) Option {
	return func(r *Review) {
		r.reviewNotifier = n
	}
}

// Database Option is specified to add the database to insert the review.
func Database(db database.Database) Option {
	return func(r *Review) {
		r.db = db
	}
}

// Logger specifies the logger used to log messages
func Logger(l log.Logger) Option {
	return func(r *Review) {
		r.log = l
	}
}

// Tracer provides the tracing
func Tracer(tp trace.Tracer) Option {
	return func(s *Review) {
		s.tracer = tp
	}
}
