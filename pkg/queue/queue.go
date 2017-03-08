package queue

import (
	"time"
)

// Queue is an interface for all queue-like impls.
type Queue interface {
	Publish(topic string, body []byte) error
	Close() error
	Subscribe(topic, channel string, handler Handler, concurrency int) error
}

// Handler is a function used for message handling.
type Handler func(msg Message) error

// Message is an abstracted interface for messages from the queue.
type Message interface {
	Finish() error
	GetBody() []byte
	GetID() []byte
	GetTimestamp() time.Time
	Requeue(delay time.Duration) error
	Touch() error
}
