package memory

import (
	"context"
	"errors"
	"net/http"

	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/protobuf/proto"
	"safer.place/internal/queue"
	"safer.place/internal/queue/otelhelper"
)

type Queue[T proto.Message] struct {
	tracer   trace.Tracer
	messages chan *Message[T]
}

type Message[T proto.Message] struct {
	// q is only used to re-queue the message when NAcked.
	q *Queue[T]

	queue.Message[T]
}

func (m *Message[T]) Body() T {
	return m.Message.Body()
}

// Ack does nothing
func (m *Message[T]) Ack() {}

// Nack restacks the message to the queue again
func (m *Message[T]) Nack() {
	m.q.messages <- m
}

// Metadata associated with the message.
func (m *Message[T]) Metadata() http.Header {
	return m.Message.Metadata()
}

// New creates a simple in memory queue based on Go channels.
func New[T proto.Message](opts ...Option[T]) *Queue[T] {
	q := &Queue[T]{
		messages: make(chan *Message[T]),
	}

	for _, opt := range opts {
		opt(q)
	}

	if err := validate(q); err != nil {
		panic(err)
	}

	return q
}

// Produce the message to the queue
func (q *Queue[T]) Produce(ctx context.Context, msg queue.Message[T]) error {
	_, span := otelhelper.StartProducerSpan(ctx, q.tracer, "incident publish")
	defer span.End()

	q.messages <- &Message[T]{Message: msg}
	span.SetStatus(codes.Ok, "")
	return nil
}

// Consume the message
func (q *Queue[T]) Consume(_ context.Context) (queue.Message[T], error) {
	msg := <-q.messages
	_, span := otelhelper.StartConsumerSpan(q.tracer, "incident receive", msg.Metadata())
	defer span.End()

	span.SetStatus(codes.Ok, "")
	return &Message[T]{
		q:       q,
		Message: msg,
	}, nil
}

var (
	errMissingTracer = errors.New("missing tracer")
)

func validate[T proto.Message](q *Queue[T]) error {
	if q.tracer == nil {
		return errMissingTracer
	}

	return nil
}
