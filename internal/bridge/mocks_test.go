package bridge

import (
	"github.com/chrisrickenbacher/lox-mqtt-bridge/internal/loxone"
	"github.com/stretchr/testify/mock"
)

// MockLoxoneProvider is a mock implementation of LoxoneProvider
type MockLoxoneProvider struct {
	mock.Mock
}

func (m *MockLoxoneProvider) Connect() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockLoxoneProvider) GetStructure() (*loxone.LoxApp3, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*loxone.LoxApp3), args.Error(1)
}

func (m *MockLoxoneProvider) EnableStatusUpdates() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockLoxoneProvider) SendCommand(cmd string) error {
	args := m.Called(cmd)
	return args.Error(0)
}

func (m *MockLoxoneProvider) GetEvents() <-chan loxone.Event {
	args := m.Called()
	return args.Get(0).(<-chan loxone.Event)
}

func (m *MockLoxoneProvider) Close() {
	m.Called()
}

// MockMQTTProvider is a mock implementation of MQTTProvider
type MockMQTTProvider struct {
	mock.Mock
}

func (m *MockMQTTProvider) Connect() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockMQTTProvider) Publish(topic string, qos byte, retained bool, payload interface{}) error {
	args := m.Called(topic, qos, retained, payload)
	return args.Error(0)
}

func (m *MockMQTTProvider) Subscribe(topic string, qos byte, callback func(topic string, payload []byte)) error {
	args := m.Called(topic, qos, callback)
	return args.Error(0)
}

func (m *MockMQTTProvider) Close() {
	m.Called()
}
