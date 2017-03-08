package queue

import (
	"sync"
	"time"

	"github.com/dchest/uniuri"
	"gopkg.in/eapache/channels.v1"
)

var defaultDelay = time.Second * 1

// Memory is an implementation of the Queue that is purely in-memory
type Memory struct {
	// topics channels chans
	queues map[string]map[string]*channels.InfiniteChannel
	mu     sync.Mutex
	stops  []chan struct{}
}

// NewMemory creates a new Memory queue
func NewMemory() (*Memory, error) {
	return &Memory{
		queues: map[string]map[string]*channels.InfiniteChannel{},
		stops:  []chan struct{}{},
	}, nil
}

// Publish sends a message to the specified topic.
func (m *Memory) Publish(topic string, body []byte) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	tm, ok := m.queues[topic]
	if !ok {
		return nil
	}

	for name, channel := range tm {
		channel.In() <- &MemoryMessage{
			Body:      body,
			ID:        []byte(uniuri.New()),
			Timestamp: time.Now(),

			Topic:   topic,
			Channel: name,
		}
	}

	return nil
}

// Close closes all queue handlers.
func (m *Memory) Close() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	for _, ch := range m.stops {
		ch <- struct{}{}
	}

	return nil
}

// Subscribe handles messages by binding to a specific topic and channel.
func (m *Memory) Subscribe(topic, channelName string, handler Handler, concurrency int) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.queues[topic]; !ok {
		m.queues[topic] = map[string]*channels.InfiniteChannel{}
	}

	topicMap := m.queues[topic]

	if _, ok := topicMap[channelName]; !ok {
		topicMap[channelName] = channels.NewInfiniteChannel()
	}

	var wg sync.WaitGroup
	wg.Add(concurrency)

	for i := 0; i < concurrency; i++ {
		// create a new stopping chan
		stop := make(chan struct{})
		m.stops = append(m.stops, stop)

		go func() {
			wg.Done()

			for {
				select {
				case msgi := <-m.queues[topic][channelName].Out():
					msg := msgi.(*MemoryMessage)

					if err := handler(msg); err != nil || msg.Requeued {
						delay := defaultDelay
						if msg.Requeued {
							delay = msg.Delay
						}

						time.AfterFunc(delay, func() {
							m.mu.Lock()
							defer m.mu.Unlock()

							m.queues[msg.Topic][msg.Channel].In() <- &MemoryMessage{
								Body:      msg.Body,
								ID:        []byte(uniuri.New()),
								Timestamp: time.Now(),

								Topic:   msg.Topic,
								Channel: msg.Channel,
							}
						})
					}
				case <-stop:
					return
				}
			}
		}()
	}

	wg.Wait()

	return nil
}

// MemoryMessage is an impl of Message
type MemoryMessage struct {
	Body      []byte
	ID        []byte
	Timestamp time.Time

	Topic   string
	Channel string

	Requeued bool
	Delay    time.Duration
}

// Finish is a nop
func (m *MemoryMessage) Finish() error {
	return nil
}

// GetBody returns body
func (m *MemoryMessage) GetBody() []byte {
	return m.Body
}

// GetID returns id
func (m *MemoryMessage) GetID() []byte {
	return m.ID
}

// GetTimestamp returns the creation timestamp
func (m *MemoryMessage) GetTimestamp() time.Time {
	return m.Timestamp
}

// Requeue requeues the message
func (m *MemoryMessage) Requeue(delay time.Duration) error {
	m.Requeued = true
	m.Delay = delay // possibly add check here for neg time cuz its in64
	return nil
}

// Touch is a nop
func (m *MemoryMessage) Touch() error {
	return nil
}
