package queue

import (
	"sync"
	"time"

	"github.com/nsqio/go-nsq"

	"github.com/oakmail/backend/pkg/config"
)

// NSQ is a queue implementation using the nsq client
type NSQ struct {
	servers   []string
	lookupds  []string
	consumers []*nsq.Consumer
	producers []*nsq.Producer
}

// NewNSQ returns a new NSQ-using queue
func NewNSQ(cfg config.NSQConfig) (*NSQ, error) {
	n := &NSQ{
		servers:   cfg.NSQdAddresses,
		lookupds:  cfg.LookupdAddresses,
		consumers: []*nsq.Consumer{},
		producers: []*nsq.Producer{},
	}

	for _, nsqd := range cfg.NSQdAddresses {
		// todo: consider using shared cfg, passing settings to newnsq etc
		producer, err := nsq.NewProducer(nsqd, nsq.NewConfig())
		if err != nil {
			return nil, err
		}
		n.producers = append(n.producers, producer)
	}

	return n, nil
}

var (
	counter uint64
	mutex   sync.Mutex
)

// Publish uses round-robin to publish a msg to one of the NSQ producers.
func (n *NSQ) Publish(topic string, body []byte) error {
	mutex.Lock()

	counter = (counter + 1) % uint64(len(n.producers))
	producer := n.producers[counter]

	mutex.Unlock()

	return producer.Publish(topic, body)
}

// Close shuts down all handlers and producers
func (n *NSQ) Close() error {
	for _, consumer := range n.consumers {
		consumer.Stop()
	}

	for _, producer := range n.producers {
		producer.Stop()
	}

	return nil
}

// Subscribe creates a new subscriber in the nsq system
func (n *NSQ) Subscribe(topic, channel string, handler Handler, concurrency int) error {
	// consider nsq config sharing
	consumer, err := nsq.NewConsumer(topic, channel, nsq.NewConfig())
	if err != nil {
		return err
	}

	consumer.AddConcurrentHandlers(nsq.HandlerFunc(func(m *nsq.Message) error {
		return handler(&NSQMessage{
			Message: m,
		})
	}), concurrency)

	if len(n.servers) > 0 {
		if err := consumer.ConnectToNSQDs(n.servers); err != nil {
			return err
		}
	}

	if len(n.lookupds) > 0 {
		if err := consumer.ConnectToNSQLookupds(n.lookupds); err != nil {
			return err
		}
	}

	n.consumers = append(n.consumers, consumer)
	return nil
}

// NSQMessage wraps nsq.Message
type NSQMessage struct {
	*nsq.Message
}

// Finish calls nsq's Finish.
func (n *NSQMessage) Finish() error {
	n.Message.Finish()
	return nil
}

// GetBody returns body of the nsq message
func (n *NSQMessage) GetBody() []byte {
	return n.Body
}

// GetID returns id of the nsq message
func (n *NSQMessage) GetID() []byte {
	return n.ID[:]
}

// GetTimestamp converts the nsq timestamp to time.Time and returns it
func (n *NSQMessage) GetTimestamp() time.Time {
	return time.Unix(0, n.Timestamp)
}

// Touch calls nsq's Touch
func (n *NSQMessage) Touch() error {
	n.Message.Touch()
	return nil
}

// Requeue calls nsq's Requeue
func (n *NSQMessage) Requeue(delay time.Duration) error {
	n.Message.Requeue(delay)
	return nil
}
