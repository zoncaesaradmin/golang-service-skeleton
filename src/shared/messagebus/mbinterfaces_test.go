package messagebus

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test Message struct creation and validation
func TestMessage_Creation(t *testing.T) {
	message := Message{
		Topic:   "test-topic",
		Key:     "test-key",
		Value:   []byte("test-value"),
		Headers: map[string]string{"header1": "value1"},
	}
	
	assert.Equal(t, "test-topic", message.Topic)
	assert.Equal(t, "test-key", message.Key)
	assert.Equal(t, []byte("test-value"), message.Value)
	assert.Equal(t, "value1", message.Headers["header1"])
}

// Test Message header manipulation
func TestMessage_HeaderManipulation(t *testing.T) {
	message := Message{
		Topic:   "test-topic",
		Value:   []byte("test-value"),
		Headers: make(map[string]string),
	}
	
	// Test adding headers
	message.Headers["header1"] = "value1"
	message.Headers["header2"] = "value2"
	
	assert.Equal(t, "value1", message.Headers["header1"])
	assert.Equal(t, "value2", message.Headers["header2"])
	assert.Len(t, message.Headers, 2)
}

// Interface compliance test
func TestInterfaces(t *testing.T) {
	// Test that we can create instances that implement the interfaces
	producer := NewProducer()
	assert.NotNil(t, producer)
	assert.Implements(t, (*Producer)(nil), producer)
	
	consumer := NewConsumer(producer)
	assert.NotNil(t, consumer)
	assert.Implements(t, (*Consumer)(nil), consumer)
}
