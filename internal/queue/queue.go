package queue

import (
	"context"

	"google.golang.org/protobuf/proto"
)

// Producer allows to publish messages
type Producer[T proto.Message] interface {
	Produce(context.Context, T) error
}

type Message[T proto.Message] interface {
	Body() T
	Ack()
	Nack()
}

// Consumer allows to consume messages
type Consumer[T proto.Message] interface {
	Consume(context.Context) (Message[T], error)
}

// Queue is a combined interface which exposes a queue which can both consume
// and produce.
type Queue[T proto.Message] interface {
	Producer[T]
	Consumer[T]
}
