package bridge

import (
	"github.com/chrisrickenbacher/lox-mqtt-bridge/internal/loxone"
)

// LoxoneProvider defines the interface required from the Loxone client.
type LoxoneProvider interface {
	Connect() error
	GetStructure() (*loxone.LoxApp3, error)
	EnableStatusUpdates() error
	SendCommand(cmd string) error
	GetEvents() <-chan loxone.Event
	Close()
}

// MQTTProvider defines the interface required from the MQTT client.
type MQTTProvider interface {
	Connect() error
	Publish(topic string, qos byte, retained bool, payload interface{}) error
	Subscribe(topic string, qos byte, callback func(topic string, payload []byte)) error
	Close()
}
