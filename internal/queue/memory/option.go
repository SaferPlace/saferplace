package memory

import (
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/protobuf/proto"
)

type Option[T proto.Message] func(*Queue[T])

// Trace provides the tracing
func Tracer[T proto.Message](tp trace.Tracer) Option[T] {
	return func(q *Queue[T]) {
		q.tracer = tp
	}
}
