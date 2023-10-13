package queue

import (
	"context"

	"google.golang.org/protobuf/proto"
)

// Producer allows to publish messages. The metadata is extracted from the producer in a way that the producer desires
type Producer[T proto.Message] interface {
	Produce(context.Context, Message[T]) error
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
