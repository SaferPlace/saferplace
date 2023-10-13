package queue

import (
	"net/http"

	"google.golang.org/protobuf/proto"
)

// Message Degines which calls must be specified
type Message[T proto.Message] interface {
	Body() T
	Ack()
	Nack()
	Metadata() http.Header
}

type message[T proto.Message] struct {
	body T
	md   http.Header
}

// NewMesage creates a new message
func NewMessage[T proto.Message](body T, md http.Header) *message[T] {
	return &message[T]{body: body, md: md}
}

func (m *message[T]) Body() T {
	return m.body
}

func (m *message[_]) Ack() {
	panic("not implemented")
}

func (m *message[_]) Nack() {
	panic("not implemented")
}

func (m *message[_]) Metadata() http.Header {
	return m.md
}
