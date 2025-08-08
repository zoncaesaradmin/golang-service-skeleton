//go:build local
// +build local

package messagebus

import (
	"context"
	"sync"
	"time"
)

// LocalProducer mock implementation for development
type LocalProducer struct {
	topics map[string][]Message
	mutex  sync.RWMutex
}

// NewProducer creates a new local producer
func NewProducer() Producer {
	return &LocalProducer{
		topics: make(map[string][]Message),
	}
}

// Send sends a message to local storage
func (p *LocalProducer) Send(ctx context.Context, message *Message) (int32, int64, error) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	message.Timestamp = time.Now()
	message.Partition = 0 // Single partition for local

	if _, exists := p.topics[message.Topic]; !exists {
		p.topics[message.Topic] = make([]Message, 0)
	}

	message.Offset = int64(len(p.topics[message.Topic]))
	p.topics[message.Topic] = append(p.topics[message.Topic], *message)

	return message.Partition, message.Offset, nil
}

// Close closes the local producer
func (p *LocalProducer) Close() error {
	return nil
}

// LocalConsumer mock implementation for development
type LocalConsumer struct {
	topics   []string
	producer *LocalProducer
	lastRead map[string]int64
	mutex    sync.RWMutex
}

// NewConsumer creates a new local consumer
func NewConsumer(producer Producer) Consumer {
	localProducer, ok := producer.(*LocalProducer)
	if !ok {
		localProducer = &LocalProducer{topics: make(map[string][]Message)}
	}

	return &LocalConsumer{
		producer: localProducer,
		lastRead: make(map[string]int64),
	}
}

// Subscribe subscribes to topics
func (c *LocalConsumer) Subscribe(topics []string) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.topics = topics
	for _, topic := range topics {
		if _, exists := c.lastRead[topic]; !exists {
			c.lastRead[topic] = -1
		}
	}

	return nil
}

// Poll polls for next available message
func (c *LocalConsumer) Poll(timeout time.Duration) (*Message, error) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	for _, topic := range c.topics {
		c.producer.mutex.RLock()
		messages, exists := c.producer.topics[topic]
		c.producer.mutex.RUnlock()

		if !exists {
			continue
		}

		nextOffset := c.lastRead[topic] + 1
		if int(nextOffset) < len(messages) {
			message := messages[nextOffset]
			c.lastRead[topic] = nextOffset
			return &message, nil
		}
	}

	// No messages available
	return nil, nil
}

// Commit commits the offset (no-op for local)
func (c *LocalConsumer) Commit(ctx context.Context, message *Message) error {
	// Local implementation doesn't need actual commit
	return nil
}

// Close closes the local consumer
func (c *LocalConsumer) Close() error {
	return nil
}
