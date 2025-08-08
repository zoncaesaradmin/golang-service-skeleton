package messagebus

import (
	"context"
	"time"
)

// Message represents a message in the message bus
type Message struct {
	Topic     string            `json:"topic"`
	Key       string            `json:"key,omitempty"`
	Value     []byte            `json:"value"`
	Headers   map[string]string `json:"headers,omitempty"`
	Partition int32             `json:"partition,omitempty"`
	Offset    int64             `json:"offset,omitempty"`
	Timestamp time.Time         `json:"timestamp"`
}

// Producer interface for publishing messages
type Producer interface {
	// Send sends a message and returns the partition and offset
	Send(ctx context.Context, message *Message) (partition int32, offset int64, err error)

	// Close closes the producer
	Close() error
}

// Consumer interface for consuming messages
type Consumer interface {
	// Subscribe subscribes to topics
	Subscribe(topics []string) error

	// Poll polls for messages with timeout
	Poll(timeout time.Duration) (*Message, error)

	// Commit manually commits the offset for a message
	Commit(ctx context.Context, message *Message) error

	// Close closes the consumer
	Close() error
}
