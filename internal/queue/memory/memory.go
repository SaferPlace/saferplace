package memory

import (
	"context"

	"google.golang.org/protobuf/proto"
	"safer.place/internal/queue"
)

type Queue[T proto.Message] struct {
	messages chan T
}

type Message[T proto.Message] struct {
	q *Queue[T]

	body T
}

func (m *Message[T]) Body() T {
	return m.body
}

// Ack does nothing
func (m *Message[T]) Ack() {}

// Nack restacks the message to the queue again
func (m *Message[T]) Nack() {
	m.q.messages <- m.body
}

// New creates a simple in memory queue based on Go channels.
func New[T proto.Message]() *Queue[T] {
	return &Queue[T]{
		messages: make(chan T),
	}
}

// Produce the message to the queue
func (q *Queue[T]) Produce(ctx context.Context, t T) error {
	q.messages <- t
	return nil
}

// Consume the message
func (q *Queue[T]) Consume(ctx context.Context) (queue.Message[T], error) {
	return &Message[T]{
		q:    q,
		body: <-q.messages,
	}, nil
}
