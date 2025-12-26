package mqtt

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/chrisrickenbacher/lox-mqtt-bridge/internal/config"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type Client struct {
	client mqtt.Client
	cfg    config.MQTTConfig
}

func NewClient(cfg config.MQTTConfig) *Client {
	broker := fmt.Sprintf("%s://%s:%d", cfg.Protocol, cfg.Host, cfg.Port)
	if cfg.Path != "" {
		broker = fmt.Sprintf("%s%s", broker, cfg.Path)
	}

	opts := mqtt.NewClientOptions()
	opts.AddBroker(broker)
	opts.SetClientID(cfg.ClientID)

	if cfg.User != "" && cfg.Pass != "" {
		opts.SetUsername(cfg.User)
		opts.SetPassword(cfg.Pass)
	}

	opts.SetAutoReconnect(true)
	opts.SetConnectRetry(true)
	opts.SetConnectRetryInterval(5 * time.Second)

	opts.OnConnect = func(c mqtt.Client) {
		slog.Info("Connected to MQTT broker", "broker", broker)
	}
	opts.OnConnectionLost = func(c mqtt.Client, err error) {
		slog.Warn("Lost connection to MQTT broker", "error", err)
	}

	client := mqtt.NewClient(opts)
	return &Client{
		client: client,
		cfg:    cfg,
	}
}

func (c *Client) Connect() error {
	slog.Info("Connecting to MQTT broker...")
	if token := c.client.Connect(); token.Wait() && token.Error() != nil {
		return token.Error()
	}
	return nil
}

func (c *Client) Publish(topic string, qos byte, retained bool, payload interface{}) error {
	var p interface{}
	switch v := payload.(type) {
	case string:
		p = v
	case []byte:
		p = v
	case int, int32, int64:
		p = fmt.Sprintf("%d", v)
	case float32, float64:
		p = fmt.Sprintf("%f", v)
	case bool:
		p = fmt.Sprintf("%t", v)
	default:
		return fmt.Errorf("unknown payload type")
	}

	token := c.client.Publish(topic, qos, retained, p)
	token.Wait()
	return token.Error()
}

// Subscribe subscribes to a topic
func (c *Client) Subscribe(topic string, qos byte, callback func(topic string, payload []byte)) error {
	token := c.client.Subscribe(topic, qos, func(client mqtt.Client, msg mqtt.Message) {
		callback(msg.Topic(), msg.Payload())
	})
	token.Wait()
	return token.Error()
}

func (c *Client) Close() {
	if c.client.IsConnected() {
		c.client.Disconnect(250)
	}
}
