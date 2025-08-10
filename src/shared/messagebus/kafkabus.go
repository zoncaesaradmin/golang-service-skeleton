//go:build !local
// +build !local

package messagebus

import (
	"context"
	"fmt"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

// KafkaProducer Kafka implementation for production (default)
type KafkaProducer struct {
	producer *kafka.Producer
}

// NewProducer creates a new Kafka producer
func NewProducer() Producer {
	config := &kafka.ConfigMap{
		"bootstrap.servers": "localhost:9092", // Default, should be configurable
	}

	producer, err := kafka.NewProducer(config)
	if err != nil {
		panic(fmt.Sprintf("Failed to create Kafka producer: %v", err))
	}

	return &KafkaProducer{
		producer: producer,
	}
}

// Send sends a message to Kafka
func (p *KafkaProducer) Send(ctx context.Context, message *Message) (int32, int64, error) {
	message.Timestamp = time.Now()

	kafkaMessage := &kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &message.Topic,
			Partition: message.Partition,
		},
		Key:       []byte(message.Key),
		Value:     message.Value,
		Timestamp: message.Timestamp,
	}

	// Add headers
	for key, value := range message.Headers {
		kafkaMessage.Headers = append(kafkaMessage.Headers, kafka.Header{
			Key:   key,
			Value: []byte(value),
		})
	}

	deliveryChan := make(chan kafka.Event, 1)
	defer close(deliveryChan)

	err := p.producer.Produce(kafkaMessage, deliveryChan)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to produce message: %w", err)
	}

	// Wait for delivery report
	select {
	case event := <-deliveryChan:
		if msg, ok := event.(*kafka.Message); ok {
			if msg.TopicPartition.Error != nil {
				return 0, 0, fmt.Errorf("delivery failed: %w", msg.TopicPartition.Error)
			}
			return msg.TopicPartition.Partition, int64(msg.TopicPartition.Offset), nil
		}
		return 0, 0, fmt.Errorf("unexpected event type")
	case <-ctx.Done():
		return 0, 0, ctx.Err()
	}
}

// SendAsync sends a message to Kafka asynchronously
func (p *KafkaProducer) SendAsync(ctx context.Context, message *Message) <-chan SendResult {
	resultChan := make(chan SendResult, 1)

	go func() {
		defer close(resultChan)

		message.Timestamp = time.Now()

		kafkaMessage := &kafka.Message{
			TopicPartition: kafka.TopicPartition{
				Topic:     &message.Topic,
				Partition: message.Partition,
			},
			Key:       []byte(message.Key),
			Value:     message.Value,
			Timestamp: message.Timestamp,
		}

		// Add headers
		for key, value := range message.Headers {
			kafkaMessage.Headers = append(kafkaMessage.Headers, kafka.Header{
				Key:   key,
				Value: []byte(value),
			})
		}

		deliveryChan := make(chan kafka.Event, 1)
		defer close(deliveryChan)

		err := p.producer.Produce(kafkaMessage, deliveryChan)
		if err != nil {
			resultChan <- SendResult{
				Partition: 0,
				Offset:    0,
				Error:     fmt.Errorf("failed to produce message: %w", err),
			}
			return
		}

		// Wait for delivery report
		select {
		case event := <-deliveryChan:
			if msg, ok := event.(*kafka.Message); ok {
				if msg.TopicPartition.Error != nil {
					resultChan <- SendResult{
						Partition: 0,
						Offset:    0,
						Error:     fmt.Errorf("delivery failed: %w", msg.TopicPartition.Error),
					}
					return
				}
				resultChan <- SendResult{
					Partition: msg.TopicPartition.Partition,
					Offset:    int64(msg.TopicPartition.Offset),
					Error:     nil,
				}
				return
			}
			resultChan <- SendResult{
				Partition: 0,
				Offset:    0,
				Error:     fmt.Errorf("unexpected event type"),
			}
		case <-ctx.Done():
			resultChan <- SendResult{
				Partition: 0,
				Offset:    0,
				Error:     ctx.Err(),
			}
		}
	}()

	return resultChan
}

// Close closes the Kafka producer
func (p *KafkaProducer) Close() error {
	p.producer.Close()
	return nil
}

// KafkaConsumer Kafka implementation for production
type KafkaConsumer struct {
	consumer *kafka.Consumer
}

// NewConsumer creates a new Kafka consumer
func NewConsumer(producer Producer) Consumer {
	config := &kafka.ConfigMap{
		"bootstrap.servers":  "localhost:9092", // Default, should be configurable
		"group.id":           "default-group",
		"auto.offset.reset":  "earliest",
		"enable.auto.commit": false, // Manual commit
	}

	consumer, err := kafka.NewConsumer(config)
	if err != nil {
		panic(fmt.Sprintf("Failed to create Kafka consumer: %v", err))
	}

	return &KafkaConsumer{
		consumer: consumer,
	}
}

// Subscribe subscribes to topics
func (c *KafkaConsumer) Subscribe(topics []string) error {
	return c.consumer.SubscribeTopics(topics, nil)
}

// Poll polls for messages
func (c *KafkaConsumer) Poll(timeout time.Duration) (*Message, error) {
	kafkaMessage, err := c.consumer.ReadMessage(timeout)
	if err != nil {
		if kafkaErr, ok := err.(kafka.Error); ok && kafkaErr.Code() == kafka.ErrTimedOut {
			return nil, nil // Timeout is not an error
		}
		return nil, err
	}

	message := &Message{
		Topic:     *kafkaMessage.TopicPartition.Topic,
		Key:       string(kafkaMessage.Key),
		Value:     kafkaMessage.Value,
		Headers:   make(map[string]string),
		Partition: kafkaMessage.TopicPartition.Partition,
		Offset:    int64(kafkaMessage.TopicPartition.Offset),
		Timestamp: kafkaMessage.Timestamp,
	}

	// Convert headers
	for _, header := range kafkaMessage.Headers {
		message.Headers[header.Key] = string(header.Value)
	}

	return message, nil
}

// Commit manually commits the offset
func (c *KafkaConsumer) Commit(ctx context.Context, message *Message) error {
	topicPartition := kafka.TopicPartition{
		Topic:     &message.Topic,
		Partition: message.Partition,
		Offset:    kafka.Offset(message.Offset + 1),
	}

	_, err := c.consumer.CommitOffsets([]kafka.TopicPartition{topicPartition})
	return err
}

// Close closes the Kafka consumer
func (c *KafkaConsumer) Close() error {
	return c.consumer.Close()
}
